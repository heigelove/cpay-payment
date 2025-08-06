package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/heigelove/cpay-payment/internal/code"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/pkg/validation"
)

type infoRequest struct {
	OrderNo string `uri:"id" binding:"required"` // 订单ID（hashID）
}

type infoResponse struct {
	Action     string `json:"action"`      // 订单号
	MerchantID string `json:"merchant_id"` // 商户ID
	Encdata    string `json:"encdata"`     // 加密数据
	Checksum   string `json:"checksum"`    // 校验和
	PrivateKey string `json:"privatekey"`  // 私钥

}

// Info 查询订单信息
// @Summary 查询订单信息
// @Description 查询订单信息
// @Tags API.payment
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id uri:string true "订单ID的hashid"
// @Success 200 {object} infoResponse
// @Failure 400 {object} code.Failure
// @Router /api/payment/info [get]
// @Security LoginToken
func (h *handler) Info() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(infoRequest)
		res := new(infoResponse)

		// 绑定URI参数
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}

		orderInfoCacheKey := fmt.Sprintf("payin:order:wakeup:%s", req.OrderNo)
		info, err := h.cache.Get(orderInfoCacheKey)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.CacheGetError,
				"获取订单信息失败").WithError(err),
			)
			return
		}
		if info == "" {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.OrderNotFound,
				"订单信息未找到"),
			)
			return
		}

		json.Unmarshal([]byte(info), res)
		if res.Action == "" {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.OrderNotFound,
				"订单信息未找到"),
			)
			return
		}

		ctx.Payload(res)
	}
}
