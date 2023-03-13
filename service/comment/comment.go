package comment

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"talkee/core"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"
	"github.com/shopspring/decimal"
)

func New(wallet *goar.Wallet,
	comments core.CommentStore,
	cfg Config,
) *CommentService {
	return &CommentService{
		wallet:   wallet,
		comments: comments,
		cfg:      cfg,
	}
}

type Config struct {
	AppName string
}

type CommentService struct {
	wallet   *goar.Wallet
	comments core.CommentStore
	cfg      Config
}

func (s *CommentService) FindArweaveSyncList(ctx context.Context, limit uint64) ([]*core.Comment, error) {
	return s.comments.FindArweaveSyncList(ctx, limit)
}

func (s *CommentService) SyncToArweave(ctx context.Context, model *core.Comment) error {
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
	chainModel.FullName = model.Creator.FullName
	chainModel.MixinIdentityNumber = model.Creator.MixinIdentityNumber
	chainModel.MixinUserID = model.Creator.MixinUserID
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

	return s.comments.UpdateCommentTxHash(ctx, model.ID, tx.ID)
}

func (s *CommentService) FavComment(ctx context.Context, commentID, userID uint64, fav bool) (comment *core.Comment, err error) {
	comment, err = s.comments.GetComment(ctx, commentID)
	if err != nil {
		return
	}

	if err = s.comments.FavComment(ctx, comment.ID, userID, fav); err != nil {
		return
	}

	return
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
