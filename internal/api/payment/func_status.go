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

type statusRequest struct {
	OrderNo string `uri:"id" binding:"required"` // 订单ID（hashID）
}

type statusResponse struct {
	OrderNo   string `json:"order_no"`   // 订单号
	Status    string `json:"status"`     // 订单状态
	ReturnUrl string `json:"return_url"` // 同步跳转地址
}

// PayInResponse 代表代收订单表
type PayInResponse struct {
	ID               uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                                              // 主键ID
	CountryID        int        `gorm:"column:country_id;index;default:0" json:"country_id"`                                                       // 国家id
	ChannelID        int        `gorm:"column:channel_id;index;not null" json:"channel_id"`                                                        // 通道号
	MerchantNo       string     `gorm:"column:merchant_no;type:varchar(50);index;not null" json:"merchant_no"`                                     // 商户号
	OrderNo          string     `gorm:"column:order_no;type:varchar(64);uniqueIndex;not null" json:"order_no"`                                     // 系统订单号
	MerchantOrderNo  string     `gorm:"column:merchant_order_no;type:varchar(64);uniqueIndex:uk_merchant_id_no;not null" json:"merchant_order_no"` // 商户订单号
	TransactionNo    string     `gorm:"column:transaction_no;type:varchar(100);default:''" json:"transaction_no"`                                  // 业务单号
	Amount           float64    `gorm:"column:amount;type:decimal(18,4)" json:"amount"`                                                            // 下单金额
	PaymentAmount    float64    `gorm:"column:payment_amount;type:decimal(18,4);default:0.0000" json:"payment_amount"`                             // 支付金额
	Fee              float64    `gorm:"column:fee;type:decimal(18,4);default:0.0000" json:"fee"`                                                   // 手续费
	DueAmount        float64    `gorm:"column:due_amount;type:decimal(18,4);default:0.0000" json:"due_amount"`                                     // 应结金额
	Status           string     `gorm:"column:status;type:varchar(20);index;not null;default:created" json:"status"`                               // 订单状态(created,success,failed,canceled)
	MerchantStatus   string     `gorm:"column:merchant_status;type:varchar(20);default:not" json:"merchant_status"`                                // 商户通知状态(notified,not)
	SettlementStatus string     `gorm:"column:settlement_status;type:varchar(20);default:not" json:"settlement_status"`                            // 结算状态(unsettled,settled)
	Goods            string     `gorm:"column:goods;type:varchar(100)" json:"goods"`                                                               // 商品
	UserID           string     `gorm:"column:user_id;type:varchar(50)" json:"user_id"`                                                            // 付款用户id
	Name             string     `gorm:"column:name;type:varchar(50);not null" json:"name"`                                                         // 付款用户名
	Phone            string     `gorm:"column:phone;type:varchar(20);not null" json:"phone"`                                                       // 付款用户手机
	Email            string     `gorm:"column:email;type:varchar(100);not null" json:"email"`                                                      // 付款用户邮箱
	Attach           string     `gorm:"column:attach;type:varchar(255)" json:"attach"`                                                             // 其它参数
	NotifyURL        string     `gorm:"column:notify_url;type:varchar(255);not null" json:"notify_url"`                                            // 回调地址
	ReturnURL        string     `gorm:"column:return_url;type:varchar(255)" json:"return_url"`                                                     // 同步跳转地址
	Remark           string     `gorm:"column:remark;type:varchar(255)" json:"remark"`                                                             // 备注
	SyncTime         *time.Time `gorm:"column:sync_time;index" json:"sync_time"`                                                                   // 同步时间
	NotifyTime       *time.Time `gorm:"column:notify_time" json:"notify_time"`                                                                     // 下发时间
	CreateAt         time.Time  `gorm:"column:create_at;index;not null" json:"create_at"`                                                          // 下单时间
	UpdatedAt        *time.Time `gorm:"column:updated_at" json:"updated_at"`                                                                       // 更新时间
	SettlementAt     *time.Time `gorm:"column:settlement_at" json:"settlement_at"`                                                                 // 结算时间
	DeletedAt        *time.Time `gorm:"column:deleted_at;index" json:"deleted_at"`                                                                 // 删除时间
}

// Status 修改订单状态
// @Summary 修改订单状态
// @Description 修改订单状态
// @Tags API.payment
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param id uri:string true "订单ID的hashid"
// @Success 200 {object} statusResponse
// @Failure 400 {object} code.Failure
// @Router /api/payment/status [get]
// @Security LoginToken
func (h *handler) Status() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(statusRequest)
		res := new(statusResponse)
		// 绑定URI参数
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}

		upiRequestCacheKey := fmt.Sprintf("payin:upi:%s:no", req.OrderNo)
		orderNo, err := h.cache.Get(upiRequestCacheKey)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.CacheGetError,
				"获取订单信息失败").WithError(err),
			)
			return
		}

		if orderNo == "" {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.OrderNotFound,
				"订单未找到"),
			)
			return
		}

		orderStatusCacheKey := fmt.Sprintf("payin:order:%s", orderNo)
		info, err := h.cache.Get(orderStatusCacheKey)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.CacheGetError,
				"获取订单状态失败").WithError(err),
			)
			return
		}
		if info == "" {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.OrderNotFound,
				"订单状态未找到").WithError(err),
			)
			return
		}
		order := &statusResponse{}
		if err := json.Unmarshal([]byte(info), order); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.CacheGetError,
				"解析订单状态失败").WithError(err),
			)
			return
		}

		res.OrderNo = req.OrderNo
		res.Status = order.Status
		res.ReturnUrl = order.ReturnUrl

		ctx.Payload(res)
	}
}
