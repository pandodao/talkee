package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"talkee/core"
	"talkee/store"
	"talkee/store/comment"
	"talkee/store/favourite"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/shopspring/decimal"
)

func New(wallet *goar.Wallet,
	comments core.CommentStore,
	rewards core.RewardStore,
	users core.UserStore,
	cfg Config,
) *CommentService {
	return &CommentService{
		wallet:   wallet,
		comments: comments,
		rewards:  rewards,
		users:    users,
		cfg:      cfg,
	}
}

type Config struct {
	AppName string
}

type CommentService struct {
	wallet   *goar.Wallet
	comments core.CommentStore
	rewards  core.RewardStore
	users    core.UserStore
	cfg      Config
}

func (s *CommentService) FindArweaveSyncList(ctx context.Context, limit uint64) ([]*core.ArweaveSyncListItem, error) {
	return s.comments.FindArweaveSyncList(ctx, limit)
}

func (s *CommentService) SyncToArweave(ctx context.Context, model *core.ArweaveSyncListItem) error {
	arAmount, err := s.wallet.Client.GetWalletBalance(s.wallet.Signer.Address)
	if err != nil {
		return err
	}

	balance, err := decimal.NewFromString(arAmount.String())
	if err != nil {
		return err
	}

	if !balance.IsPositive() {
		return fmt.Errorf("arwallet balance not sufficient")
	}

	if model.ArweaveTxHash != "" {
		return nil
	}

	chainModel := &core.CommentChian{}
	chainModel.Content = model.Content
	chainModel.CreatedAt = model.CreatedAt
	chainModel.FullName = model.FullName
	chainModel.MixinIdentityNumber = model.MixinIdentityNumber
	chainModel.MixinUserID = model.MixinUserID
	chainModel.Slug = model.Slug

	data, err := json.Marshal(chainModel)
	if err != nil {
		return err
	}

	tx, err := s.wallet.SendData(
		data, // Data bytes
		[]types.Tag{
			types.Tag{
				Name:  "Content-Type",
				Value: "application/json",
			},
			types.Tag{
				Name:  "App-Name",
				Value: s.cfg.AppName,
			},
			types.Tag{
				Name:  "Author",
				Value: strconv.FormatUint(model.UserID, 10),
			},
		},
	)

	if err != nil {
		return err
	}

	model.ArweaveTxHash = tx.ID

	return s.comments.UpdateCommentTxHash(ctx, model.CommentID, tx.ID)
}

func (s *CommentService) FavComment(ctx context.Context, commentID, userID uint64, fav bool) (*core.Comment, error) {
	com, err := s.comments.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	if err := store.Transaction(func(tx *store.Handler) error {
		favourites := favourite.New(tx)
		comments := comment.New(tx)

		if err = comments.FavComment(ctx, com.ID, userID, fav); err != nil {
			return err
		}

		if fav {
			if err = favourites.CreateFavourite(ctx, userID, com.ID); err != nil {
				return err
			}
		} else {
			if err = favourites.DeleteFavourite(ctx, userID, com.ID); err != nil {
				return err
			}
		}
		return err
	}); err != nil {
		return nil, err
	}

	if fav {
		com.FavorCount++
	} else {
		com.FavorCount--
	}

	return com, nil
}

func (s *CommentService) CreateComment(ctx context.Context, userID, siteID uint64, slug, content string) (*core.Comment, error) {
	commentID, err := s.comments.CreateComment(ctx, userID, siteID, slug, content)
	if err != nil {
		return nil, err
	}

	comment, err := s.comments.GetComment(ctx, commentID)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *CommentService) GetCommentsWithRewards(ctx context.Context, siteID uint64, slug string, page, limit uint64, orderBy, order string) ([]*core.Comment, error) {
	offset := page * limit
	comms, err := s.comments.GetComments(ctx, siteID, slug, offset, limit, orderBy, order)
	if err != nil {
		return nil, err
	}

	if err := s.WithRewards(ctx, comms...); err != nil {
		return nil, err
	}

	if err := s.WithUsers(ctx, comms...); err != nil {
		return nil, err
	}

	return comms, nil
}

func (s *CommentService) GetCommentWithReward(ctx context.Context, id uint64) (*core.Comment, error) {
	cm, err := s.comments.GetComment(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.WithRewards(ctx, cm); err != nil {
		return nil, err
	}

	if err := s.WithUsers(ctx, cm); err != nil {
		return nil, err
	}

	return cm, nil
}

func (s *CommentService) WithRewards(ctx context.Context, comments ...*core.Comment) error {
	cmtIDs := []uint64{}
	for _, c := range comments {
		cmtIDs = append(cmtIDs, c.ID)
	}

	rewards, err := s.rewards.FindRewardsByCommentIDs(ctx, cmtIDs)
	if err != nil {
		return err
	}

	for _, r := range rewards {
		for _, c := range comments {
			if r.ObjectID == c.ID {
				c.Reward = r
			}
		}
	}

	return nil
}

func (s *CommentService) WithUsers(ctx context.Context, comments ...*core.Comment) error {
	userIDs := []uint64{}
	for _, c := range comments {
		userIDs = append(userIDs, c.UserID)
	}

	users, err := s.users.GetUserByIDs(ctx, userIDs)
	if err != nil {
		return err
	}

	for _, c := range comments {
		for _, u := range users {
			if u.ID == c.UserID {
				c.Creator = u
			}
		}
	}

	return nil
}
