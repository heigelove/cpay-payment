package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/repository/redis"

	"go.uber.org/zap"
)

type handler struct {
	logger *zap.Logger
	cache  redis.Repo
}

func New(logger *zap.Logger, cache redis.Repo) *handler {
	return &handler{
		logger: logger,
		cache:  cache,
	}
}

// Index 首页
func (h *handler) Index() core.HandlerFunc {
	return func(ctx core.Context) {
		ctx.HTML("index", nil)
	}
}

// QRCode 扫码支付
func (h *handler) QRCode() core.HandlerFunc {
	return func(ctx core.Context) {
		ctx.HTML("qrcode", nil)
	}
}

// Redirect 跳转支付
func (h *handler) Redirect() core.HandlerFunc {
	return func(ctx core.Context) {
		orderId := ctx.QueryValue("order_no")
		payData, err := h.cache.Get(fmt.Sprintf("payin:redirect:%s:data", orderId))
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, http.StatusInternalServerError, "get pay data error"))
			return
		}
		payDataMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(payData), &payDataMap)
		if err != nil {
			ctx.AbortWithError(core.Error(http.StatusInternalServerError, http.StatusInternalServerError, "unmarshal pay data error"))
			return
		}

		if pkgName, ok := payDataMap["pkg_name"].(string); ok && pkgName != "" {
			raw, ok := payDataMap["raw"].(string)
			if !ok {
				ctx.AbortWithError(core.Error(http.StatusInternalServerError, http.StatusInternalServerError, "get raw data error"))
				return
			}
			rawMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(raw), &rawMap)
			if err != nil {
				ctx.AbortWithError(core.Error(http.StatusInternalServerError, http.StatusInternalServerError, "unmarshal raw data error"))
				return
			}
			ctx.HTML("redirect_"+pkgName, rawMap)
		} else {
			ctx.HTML("redirect_", nil)
		}
	}
}
