package tip

import (
	"net/http"
	"talkee/core"
	"talkee/handler/render"
	"talkee/session"

	"github.com/fox-one/pkg/httputil/param"
	"github.com/fox-one/pkg/uuid"
	"github.com/shopspring/decimal"
)

type CreateTipPayload struct {
	SiteID         uint64                 `json:"site_id"`
	Slug           string                 `json:"slug"`
	AirdropType    string                 `json:"airdrop_type"`
	StrategyName   string                 `json:"strategy_name"`
	StrategyParams map[string]interface{} `json:"strategy_params"`
	AssetID        string                 `json:"asset_id"`
	Amount         decimal.Decimal        `json:"amount"`
	Memo           string                 `json:"memo"`
	RedirectURL    string                 `json:"redirect_url"`
}

func CreateTip(tipz core.TipService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user, found := session.UserFrom(ctx)
		if !found {
			render.Error(w, http.StatusUnauthorized, nil)
			return
		}

		body := &CreateTipPayload{}
		if err := param.Binding(r, body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		tip, err := tipz.CreateTip(ctx, &core.Tip{
			UUID:           uuid.New(),
			UserID:         user.ID,
			SiteID:         body.SiteID,
			Slug:           body.Slug,
			AirdropType:    body.AirdropType,
			StrategyName:   body.StrategyName,
			StrategyParams: body.StrategyParams,
			AssetID:        body.AssetID,
			Amount:         body.Amount,
			Memo:           body.Memo,
		}, body.RedirectURL)
		if err != nil {
			render.Error(w, http.StatusInternalServerError, err)
			return
		}

		render.JSON(w, tip)
	}
}
