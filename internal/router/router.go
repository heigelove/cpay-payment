package router

import (
	"github.com/heigelove/cpay-payment/internal/alert"
	"github.com/heigelove/cpay-payment/internal/metrics"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/repository/redis"
	"github.com/heigelove/cpay-payment/pkg/errors"

	"go.uber.org/zap"
)

type resource struct {
	mux    core.Mux
	logger *zap.Logger
	cache  redis.Repo
}

type Server struct {
	Mux   core.Mux
	Cache redis.Repo
}

func NewHTTPServer(logger *zap.Logger, cronLogger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	// 初始化 Cache
	cacheRepo, err := redis.New()
	if err != nil {
		logger.Fatal("new cache err", zap.Error(err))
	}
	r.cache = cacheRepo

	mux, err := core.New(logger,
		core.WithEnableCors(),
		core.WithEnableRate(),
		core.WithEnableSafeCheck(),
		core.WithAlertNotify(alert.NotifyHandler(logger)),
		core.WithRecordMetrics(metrics.RecordHandler(logger)),
	)

	if err != nil {
		panic(err)
	}

	r.mux = mux

	// 设置 Render 路由
	setRenderRouter(r)

	// 设置 API 路由
	setApiRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Cache = r.cache

	return s, nil
}
