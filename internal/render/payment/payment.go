package payment

import (
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
		ctx.HTML("redirect", nil)
	}
}
