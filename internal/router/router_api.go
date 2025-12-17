package router

import (
	"github.com/heigelove/cpay-payment/internal/api/payment"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
)

func setApiRouter(r *resource) {
	// 需要签名验证，无需登录验证，无需 RBAC 权限验证
	login := r.mux.Group("/api")
	{
		paymentHandler := payment.New(r.logger, r.cache)
		login.GET("/payment/status/:id", core.AliasForRecordMetrics("/api/payment/status"), paymentHandler.Status())
		login.GET("/payment/info/:id", core.AliasForRecordMetrics("/api/payment/info"), paymentHandler.Info())
		login.GET("/payment/upi/:id", core.AliasForRecordMetrics("/api/payment/upi"), paymentHandler.Upi())
		login.GET("/payment/stats", core.AliasForRecordMetrics("/api/payment/stats"), paymentHandler.Stats())
	}
}
