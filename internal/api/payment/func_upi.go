package payment

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/heigelove/cpay-payment/internal/code"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/pkg/validation"
	"go.uber.org/zap"
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

type WalletInfo struct {
	WalletName string `json:"wallet_name"` // 钱包名称
	WalletCode string `json:"wallet_code"` // 钱包代码
	PayUrl     string `json:"pay_url"`     // 钱包链接
}

type PayinOrderDao struct {
	ID               uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                                              // 主键ID
	CountryID        int        `gorm:"column:country_id;index;default:0" json:"country_id"`                                                       // 国家id
	ChannelID        int        `gorm:"column:channel_id;index;not null" json:"channel_id"`                                                        // 通道号
	MerchantNo       string     `gorm:"column:merchant_no;type:varchar(50);index;not null" json:"merchant_no"`                                     // 商户号
	OrderNo          string     `gorm:"column:order_no;type:varchar(64);uniqueIndex;not null" json:"order_no"`                                     // 系统订单号
	MerchantOrderNo  string     `gorm:"column:merchant_order_no;type:varchar(64);uniqueIndex:uk_merchant_id_no;not null" json:"merchant_order_no"` // 商户订单号
	TransactionNo    string     `gorm:"column:transaction_no;type:varchar(100);default:''" json:"transaction_no"`                                  // 业务单号
	UTR              string     `gorm:"column:utr;type:varchar(50);default:''" json:"utr"`                                                         // 银行交易参考号
	Amount           float64    `gorm:"column:amount;type:decimal(18,4)" json:"amount"`                                                            // 下单金额
	PaymentAmount    float64    `gorm:"column:payment_amount;type:decimal(18,4);default:0.0000" json:"payment_amount"`                             // 支付金额
	Fee              float64    `gorm:"column:fee;type:decimal(18,4);default:0.0000" json:"fee"`                                                   // 手续费
	DueAmount        float64    `gorm:"column:due_amount;type:decimal(18,4);default:0.0000" json:"due_amount"`                                     // 应结金额
	Profit           float64    `gorm:"column:profit;type:decimal(18,4);default:0.0000" json:"profit"`                                             // 利润
	Status           string     `gorm:"column:status;type:varchar(20);index;not null;default:created" json:"status"`                               // 订单状态(created,success,failed,canceled)
	MerchantStatus   string     `gorm:"column:merchant_status;type:varchar(20);default:not" json:"merchant_status"`                                // 商户通知状态(notified,not)
	SettlementStatus string     `gorm:"column:settlement_status;type:varchar(20);default:unsettled" json:"settlement_status"`                      // 结算状态(unsettled,settled)
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
				"Failed to obtain the order URL").WithError(err),
			)
			return
		}
		if url == "" {
			ctx.AbortWithError(core.Error(
				http.StatusNotFound,
				code.OrderNotFound,
				"Failed to obtain the order URL"),
			)
			return
		}

		go func(logger *zap.Logger) {
			err := h.statsCount(req)
			if err != nil {
				logger.Error("Failed to record UPI order stats", zap.Error(err))
			}
		}(ctx.Logger())

		var walletList []WalletInfo
		err = json.Unmarshal([]byte(url), &walletList)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusInternalServerError,
				code.OrderNotFound,
				"Failed to parse the order URL").WithError(err),
			)
			return
		}
		var upiUrl string
		data := upiResponseData{}
		for _, wallet := range walletList {
			if wallet.WalletCode == "upi" {
				upiUrl = wallet.PayUrl
			}
			if wallet.WalletCode == "bhim" {
				data.Bhim = wallet.PayUrl
			}
			if wallet.WalletCode == "gpay" {
				data.Gpay = wallet.PayUrl
			}
			if wallet.WalletCode == "paytm" {
				data.Paytm = wallet.PayUrl
			}
			if wallet.WalletCode == "phonepe" {
				data.Phonepe = wallet.PayUrl
			}
		}
		data.Default = upiUrl
		data.UPI = upiUrl
		if upiUrl != "" {
			data.PA = extractPAFromUPI(upiUrl)
		}

		res.Code = "0000"
		res.Msg = "success"
		res.Data = data

		ctx.Payload(res)
	}
}

func (h *handler) statsCount(req *upiRequest) error {
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
	// key: stats:upi:access:count
	// member: channel_id
	// score: 访问次数
	statsKey := "stats:upi:access:count"
	channelMember := fmt.Sprintf("channel_%d", payInOrder.ChannelID)

	// 增加通道访问计数
	_, err = h.cache.ZIncrBy(statsKey, 1.0, channelMember)
	if err != nil {
		return fmt.Errorf("Failed to increment channel stats: %v", err)
	}

	// 同时记录按日期分组的统计
	today := time.Now().Format("2006-01-02")
	dailyStatsKey := fmt.Sprintf("stats:upi:daily:%s", today)
	_, err = h.cache.ZIncrBy(dailyStatsKey, 1.0, channelMember)
	if err != nil {
		return fmt.Errorf("Failed to increment daily channel stats: %v", err)
	}

	// 设置每日统计的过期时间为7天
	h.cache.Expire(dailyStatsKey, 7*24*time.Hour)

	return nil
}
