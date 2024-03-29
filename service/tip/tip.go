package tip

import (
	"context"
	"fmt"
	"sort"
	"talkee/core"
	"talkee/internal/mixpay"
	"talkee/store"
	"talkee/store/reward"
	"talkee/store/tip"
	"talkee/util"

	"github.com/fox-one/pkg/uuid"
	"github.com/shopspring/decimal"
)

func New(
	cfg Config,
	mixpayClient *mixpay.Client,
	tips core.TipStore,
	comments core.CommentStore,
	rewards core.RewardStore,
	users core.UserStore,
	sites core.SiteStore,
	commentz core.CommentService,
) *service {
	return &service{
		cfg:          cfg,
		mixpayClient: mixpayClient,
		tips:         tips,
		comments:     comments,
		rewards:      rewards,
		users:        users,
		sites:        sites,
		commentz:     commentz,
	}
}

type Config struct {
	ClientID          string
	MixpayPayeeID     string
	MixpayCallbackURL string
}

type service struct {
	cfg          Config
	mixpayClient *mixpay.Client
	tips         core.TipStore
	comments     core.CommentStore
	rewards      core.RewardStore
	users        core.UserStore
	sites        core.SiteStore
	commentz     core.CommentService
}

func (s *service) CreateTip(ctx context.Context, tp *core.Tip, redirectURL string) (*core.Tip, error) {
	if err := tp.Validate(); err != nil {
		return nil, err
	}

	tp.Status = core.TipStatusInit
	tipID, err := s.tips.CreateTip(ctx, tp)
	if err != nil {
		return nil, err
	}

	tp, err = s.tips.GetTip(ctx, tipID)
	if err != nil {
		return nil, err
	}

	req := mixpay.CreateOneTimePaymentRequest{
		PayeeId:           s.cfg.MixpayPayeeID,
		QuoteAssetId:      tp.AssetID,
		QuoteAmount:       tp.Amount.StringFixed(8),
		SettlementAssetId: tp.AssetID,
		OrderId:           tp.UUID,
		TraceId:           tp.UUID,
		CallbackUrl:       s.cfg.MixpayCallbackURL,
		ReturnTo:          redirectURL,
		FailedReturnTo:    redirectURL,
	}

	resp, err := s.mixpayClient.CreateOneTimePayment(ctx, req)

	if err != nil {
		return nil, err
	}

	tp.MixpayCode = resp.Code

	return tp, nil
}

func (s *service) FillTipByMixpay(ctx context.Context, tipUUID string) (*core.Tip, error) {
	tp, err := s.tips.GetTipByUUID(ctx, tipUUID)
	if err != nil {
		return nil, err
	}

	if tp.Status != core.TipStatusInit {
		return nil, core.ErrIncorrectTipStatus
	}

	resp, err := s.mixpayClient.GetPaymentResult(ctx, mixpay.GetPaymentResultRequest{
		OrderId: tipUUID,
		TraceId: tipUUID,
		PayeeId: s.cfg.ClientID,
	})
	if err != nil {
		return nil, err
	}

	fmt.Printf("resp: %+v\n", resp)
	if resp.PayeeID != s.cfg.MixpayPayeeID || resp.QuoteAssetID != tp.AssetID {
		return nil, core.ErrInvalidMixpayPaymentInfo
	}

	amount, err := decimal.NewFromString(resp.QuoteAmount)
	if err != nil {
		return nil, err
	}

	// @TODO not a good way to check
	if amount.LessThan(tp.Amount.Sub(decimal.NewFromInt(1))) {
		return nil, core.ErrInvalidMixpayPaymentInfo
	}

	if resp.Status == "success" {
		if err := s.tips.UpdateTipStatus(ctx, tp.ID, core.TipStatusFilled); err != nil {
			return nil, err
		}
	} else if resp.Status == "failed" {
		if err := s.tips.UpdateTipStatus(ctx, tp.ID, core.TipStatusFailed); err != nil {
			return nil, err
		}
	}

	return tp, nil
}

func (s *service) GetTipsBySlug(ctx context.Context, siteID uint64, slug string, limit int) ([]*core.Tip, error) {
	tps, err := s.tips.GetTipsBySlug(ctx, siteID, slug, limit)
	if err != nil {
		return nil, err
	}

	if err := s.WithUsers(ctx, tps...); err != nil {
		return nil, err
	}

	return tps, nil
}

func (s *service) WithUsers(ctx context.Context, tps ...*core.Tip) error {
	userIDs := []uint64{}
	for _, tp := range tps {
		userIDs = append(userIDs, tp.UserID)
	}

	users, err := s.users.GetUserByIDs(ctx, userIDs)
	if err != nil {
		return err
	}

	for _, tp := range tps {
		for _, u := range users {
			if u.ID == tp.UserID {
				tp.User = u
			}
		}
	}

	return nil
}

func (s *service) ProcessPendingTip(ctx context.Context, tip *core.Tip) error {
	rs, err := s.rewards.GetRewardsByTipIDAndStatus(ctx, tip.ID, core.RewardStatusCreated)
	if err != nil {
		return err
	}

	if len(rs) == 0 {
		if err := s.tips.UpdateTipStatus(ctx, tip.ID, core.TipStatusDelivered); err != nil {
			return err
		}
	}

	return nil
}

func (s *service) ProcessFilledTip(ctx context.Context, tip *core.Tip) error {
	switch tip.AirdropType {
	case core.AirdropTypeComments:
		return s.ProcessCommentsType(ctx, tip)
	case core.AirdropTypeSlug:
		return s.ProcessSlugType(ctx, tip)
	case core.AirdropTypeComment:
		return s.ProcessCommentType(ctx, tip)
	case core.AirdropTypeUser:
		return s.ProcessUserType(ctx, tip)
	}

	return core.ErrUnsupportedTipType
}

func (s *service) ProcessSlugType(ctx context.Context, tp *core.Tip) error {
	if tp.SiteID != 0 {
		site, err := s.sites.GetSite(ctx, tp.SiteID)
		if err != nil {
			return err
		}

		owner, err := s.users.GetUser(ctx, site.UserID)
		if err != nil {
			return err
		}

		if err := store.Transaction(func(tx *store.Handler) error {
			rewards := reward.New(tx)
			tips := tip.New(tx)

			rw := &core.Reward{
				ObjectType:  core.RewardObjectTypeSlug,
				ObjectID:    site.ID,
				RecipientID: owner.MixinUserID,
				Amount:      tp.Amount,
				TipID:       tp.ID,
				Memo:        tp.Memo,
				SiteID:      tp.SiteID,
				AssetID:     tp.AssetID,
				Status:      core.RewardStatusCreated,
			}

			rw.TraceID = uuid.Modify(s.cfg.ClientID, fmt.Sprintf("reward-generated-by(%d, %s, %d)", tp.ID, rw.ObjectType, rw.ObjectID))
			if err := rewards.CreateReward(ctx, rw); err != nil {
				return err
			}

			if err := tips.UpdateTipStatus(ctx, tp.ID, core.TipStatusPending); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ProcessUserType(ctx context.Context, tp *core.Tip) error {
	if tp.OpponentID != 0 {
		user, err := s.users.GetUser(ctx, tp.OpponentID)
		if err != nil {
			return err
		}

		if err := store.Transaction(func(tx *store.Handler) error {
			rewards := reward.New(tx)
			tips := tip.New(tx)

			rw := &core.Reward{
				ObjectType:  core.RewardObjectTypeUser,
				ObjectID:    tp.OpponentID,
				RecipientID: user.MixinUserID,
				Amount:      tp.Amount,
				TipID:       tp.ID,
				Memo:        tp.Memo,
				SiteID:      tp.SiteID,
				AssetID:     tp.AssetID,
				Status:      core.RewardStatusCreated,
			}

			rw.TraceID = uuid.Modify(s.cfg.ClientID, fmt.Sprintf("reward-generated-by(%d, %s, %d)", tp.ID, rw.ObjectType, rw.ObjectID))
			if err := rewards.CreateReward(ctx, rw); err != nil {
				return err
			}

			if err := tips.UpdateTipStatus(ctx, tp.ID, core.TipStatusPending); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ProcessCommentType(ctx context.Context, tp *core.Tip) error {
	if tp.OpponentID != 0 {
		cmt, err := s.commentz.GetCommentWithReward(ctx, tp.OpponentID)
		if err != nil {
			return err
		}

		if err := store.Transaction(func(tx *store.Handler) error {
			rewards := reward.New(tx)
			tips := tip.New(tx)

			rw := &core.Reward{
				ObjectType:  core.RewardObjectTypeComment,
				ObjectID:    cmt.ID,
				RecipientID: cmt.Creator.MixinUserID,
				Amount:      tp.Amount,
				TipID:       tp.ID,
				Memo:        tp.Memo,
				SiteID:      tp.SiteID,
				AssetID:     tp.AssetID,
				Status:      core.RewardStatusCreated,
			}

			rw.TraceID = uuid.Modify(s.cfg.ClientID, fmt.Sprintf("reward-generated-by(%d, %s, %d)", tp.ID, rw.ObjectType, rw.ObjectID))
			if err := rewards.CreateReward(ctx, rw); err != nil {
				return err
			}

			if err := tips.UpdateTipStatus(ctx, tp.ID, core.TipStatusPending); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *service) ProcessCommentsType(ctx context.Context, tp *core.Tip) error {
	params := tp.StrategyParams.ToMap()
	rws := make([]*core.Reward, 0)
	switch tp.StrategyName {
	case core.TipStrategyAvg:
		{
			cmts, err := s.comments.GetAllCommentsBySiteSlug(ctx, tp.SiteID, tp.Slug)
			if err != nil {
				return err
			}

			if err := s.commentz.WithUsers(ctx, cmts...); err != nil {
				return err
			}

			shares := len(cmts)
			if shares == 0 {
				return nil
			}

			amount := tp.Amount.Div(decimal.NewFromInt(int64(shares)))
			if amount.LessThan(util.OneSatoshi) {
				return nil
			}

			for _, cmt := range cmts {
				rws = append(rws, &core.Reward{
					ObjectType:  core.RewardObjectTypeComment,
					ObjectID:    cmt.ID,
					RecipientID: cmt.Creator.MixinUserID,
					Amount:      amount,
				})
			}
		}
	case core.TipStrategyTopN:
		{
			var topN int
			topNAny, found := params["n"]
			if !found {
				topN = 3
			} else {
				topN = int(topNAny.(float64))
			}

			cmts, err := s.comments.GetAllCommentsBySiteSlug(ctx, tp.SiteID, tp.Slug)
			if err != nil {
				return err
			}

			if len(cmts) == 0 {
				return nil
			}

			if len(cmts) < topN {
				topN = len(cmts)
			}

			sort.SliceStable(cmts, func(i, j int) bool {
				return cmts[i].FavorCount > cmts[j].FavorCount
			})

			// find top n, if n + 1 has same favor count with n, then n + 1 is also in top n
			topNCmts := make([]*core.Comment, 0)
			for i, cmt := range cmts {
				if i < topN {
					topNCmts = append(topNCmts, cmt)
				} else {
					if cmt.FavorCount == cmts[i-1].FavorCount {
						topNCmts = append(topNCmts, cmt)
					} else {
						break
					}
				}
			}

			if err := s.commentz.WithUsers(ctx, topNCmts...); err != nil {
				return err
			}

			amount := tp.Amount.Div(decimal.NewFromInt(int64(topN)))
			for _, cmt := range topNCmts {
				rws = append(rws, &core.Reward{
					ObjectType:  core.RewardObjectTypeComment,
					ObjectID:    cmt.ID,
					RecipientID: cmt.Creator.MixinUserID,
					Amount:      amount,
				})
			}
		}
	default:
		return core.ErrUnsupportedTipStrategy
	}

	if err := store.Transaction(func(tx *store.Handler) error {
		rewards := reward.New(tx)
		tips := tip.New(tx)

		for _, rw := range rws {
			rw.TipID = tp.ID
			rw.Memo = tp.Memo
			rw.SiteID = tp.SiteID
			rw.AssetID = tp.AssetID
			rw.Status = core.RewardStatusCreated
			rw.TraceID = uuid.Modify(s.cfg.ClientID, fmt.Sprintf("reward-generated-by(%d, %s, %d)", tp.ID, rw.ObjectType, rw.ObjectID))
			err := rewards.CreateReward(ctx, rw)
			if err != nil {
				return err
			}
		}

		if err := tips.UpdateTipStatus(ctx, tp.ID, core.TipStatusPending); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
