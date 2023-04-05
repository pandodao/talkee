package payment

import (
	"net/http"

	"talkee/core"
	"talkee/handler/render"

	"github.com/fox-one/pkg/httputil/param"
	"github.com/fox-one/pkg/logger"
)

func HandleMixpayWebhook(tipz core.TipService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx).WithField("handler", "mixpay_callback")
		var body struct {
			OrderID string `json:"orderId"`
			TraceID string `json:"traceId"`
			PayeeID string `json:"payeeId"`
		}

		if err := param.Binding(r, &body); err != nil {
			render.Error(w, http.StatusBadRequest, err)
			return
		}

		if _, err := tipz.FillTipByMixpay(ctx, body.OrderID); err != nil {
			log.WithError(err).Println("fill tip by mixpay failed")
		}

		render.JSON(w, map[string]any{"code": "0"})
	}
}
