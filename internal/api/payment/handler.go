package payment

import (
	"github.com/heigelove/cpay-payment/configs"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/repository/redis"
	"github.com/heigelove/cpay-payment/pkg/hash"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Status 查询订单状态
	// @Tags API.payment
	// @Router /api/payment/status/:id [get]
	Status() core.HandlerFunc

	// Info 查询订单状态
	// @Tags API.payment
	// @Router /api/payment/info/:id [get]
	Info() core.HandlerFunc

	// Upi 查询订单状态
	// @Tags API.payment
	// @Router /api/payment/upi/:id [get]
	Upi() core.HandlerFunc

	// Stats 记录订单点击统计
	// @Tags API.payment
	// @Router /api/payment/stats/:id [get]
	Stats() core.HandlerFunc
}

type handler struct {
	logger  *zap.Logger
	cache   redis.Repo
	hashids hash.Hash
}

func New(logger *zap.Logger, cache redis.Repo) Handler {
	return &handler{
		logger:  logger,
		cache:   cache,
		hashids: hash.New(configs.Get().HashIds.Secret, configs.Get().HashIds.Length),
	}
}

func (h *handler) i() {}
