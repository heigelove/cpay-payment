package router

import (
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/render/payment"
)

func setRenderRouter(r *resource) {
	renderPayment := payment.New(r.logger, r.cache)

	// 无需记录日志，无需 RBAC 权限验证
	notRBAC := r.mux.Group("", core.DisableTraceLog, core.DisableRecordMetrics)
	{
		notRBAC.GET("/payment/qrcode", renderPayment.QRCode())
		notRBAC.GET("/payment/redirect", renderPayment.Redirect())
		notRBAC.GET("/payment/index", renderPayment.Index())
	}
}
