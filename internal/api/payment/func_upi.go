package payment

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/heigelove/cpay-payment/internal/code"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/pkg/validation"
)

type upiRequest struct {
	OrderNo string `uri:"id" binding:"required"` // 订单ID（hashID）
}

type upiResponse struct {
	Code string          `json:"code"`
	Msg  string          `json:"msg"`
	Data upiResponseData `json:"data"`
}

type upiResponseData struct {
	Default string `json:"default"`
	Bhim    string `json:"bhim"`
	Gpay    string `json:"gpay"`
	Paytm   string `json:"paytm"`
	Phonepe string `json:"phonepe"`
	UPI     string `json:"upi"`
	PA      string `json:"pa"`
}

// extractPAFromUPI 从UPI字符串中提取pa参数的值
// 使用正则表达式查找 pa= 到 &pn= 之间的内容
func extractPAFromUPI(upiString string) string {
	// 找到 pa= 到 & 之间的部分
	re := regexp.MustCompile(`pa=([^&]+)`)
	match := re.FindStringSubmatch(upiString)
	if len(match) >= 2 {
		return match[1]
	}
	return ""
}

// Upi 查询订单状态
// @Summary 查询订单状态
// @Description 查询订单状态
// @Tags API.payment
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id uri:string true "订单ID的hashid"
// @Success 200 {object} upiResponse
// @Failure 400 {object} code.Failure
// @Router /api/payment/upi [get]
// @Security LoginToken
func (h *handler) Upi() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(upiRequest)
		res := new(upiResponse)

		// 绑定URI参数
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}

		orderStatusCacheKey := fmt.Sprintf("payin:upi:%s:url", req.OrderNo)
		url, err := h.cache.Get(orderStatusCacheKey)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.CacheGetError,
				"获取订单URL失败").WithError(err),
			)
			return
		}
		if url == "" {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.OrderNotFound,
				"获取订单URL失败"),
			)
			return
		}
		baseUrl := strings.Replace(url, "upi://", "", 1)

		data := upiResponseData{}
		data.Default = url
		data.Bhim = fmt.Sprintf("bhim://upi/%s", baseUrl)
		data.Gpay = fmt.Sprintf("gpay://upi/%s", baseUrl)
		data.Paytm = fmt.Sprintf("paytmmp://%s", baseUrl)
		data.Phonepe = fmt.Sprintf("phonepe://%s", baseUrl)
		data.UPI = url
		data.PA = extractPAFromUPI(url)

		res.Code = "0000"
		res.Msg = "success"
		res.Data = data

		ctx.Payload(res)
	}
}
