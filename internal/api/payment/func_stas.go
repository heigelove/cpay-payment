package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/heigelove/cpay-payment/internal/code"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/pkg/validation"
)

type statsRequest struct {
	OrderNo string `form:"id" binding:"required"`       // 订单ID
	PayType string `form:"pay_type" binding:"required"` // 支付类型
}

func (h *handler) Stats() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(statsRequest)
		// 绑定URI参数
		if err := ctx.ShouldBindForm(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}
		h.statsClick(req)
		ctx.Payload(nil)
	}
}

func (h *handler) statsClick(req *statsRequest) error {
	if req.PayType != "paytm" && req.PayType != "phonepe" && req.PayType != "gpay" && req.PayType != "bhim" {
		return fmt.Errorf("Unsupported payment type for stats")
	}
	upiRequestCacheKey := fmt.Sprintf("payin:upi:%s:no", req.OrderNo)
	orderNo, err := h.cache.Get(upiRequestCacheKey)
	if err != nil {
		return err
	}

	if orderNo == "" {
		return fmt.Errorf("Order not found")
	}

	orderCacheKey := fmt.Sprintf("payin:order:%s", orderNo)
	var payInOrder PayinOrderDao
	orderJson, err := h.cache.Get(orderCacheKey)
	if err != nil {
		return err
	}

	if orderJson == "" {
		return fmt.Errorf("Order not found")
	}

	if err := json.Unmarshal([]byte(orderJson), &payInOrder); err != nil {
		return err
	}

	// 使用有序集合记录每个通道的访问次数
	// key: stats:upi:payment:click
	// member: channel_id
	// score: 访问次数
	statsKey := "stats:upi:payment:" + req.PayType + ":click"
	channelMember := fmt.Sprintf("channel_%d", payInOrder.ChannelID)

	// 增加通道访问计数
	_, err = h.cache.ZIncrBy(statsKey, 1.0, channelMember)
	if err != nil {
		return fmt.Errorf("Failed to increment channel stats: %v", err)
	}

	// 同时记录按日期分组的统计
	today := time.Now().Format("2006-01-02")
	dailyStatsKey := fmt.Sprintf("stats:upi:payment:daily:%s", today)
	_, err = h.cache.ZIncrBy(dailyStatsKey, 1.0, channelMember)
	if err != nil {
		return fmt.Errorf("Failed to increment daily channel stats: %v", err)
	}

	// 设置每日统计的过期时间为7天
	h.cache.Expire(dailyStatsKey, 7*24*time.Hour)

	return nil
}
