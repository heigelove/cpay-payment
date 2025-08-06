package query

import (
	"encoding/json"
	"sync"

	"github.com/heigelove/cpay-payment/internal/repository/redis"
	"go.uber.org/zap"
)

// 队列名称常量
const (
	OrderQueryQueue = "order_query_tasks" // 订单查询任务队列
)

// OrderQueryTask 订单查询任务结构
type OrderQueryTask struct {
	OrderNo       string `json:"order_no"`       // 系统订单号
	TransactionNo string `json:"transaction_no"` // 商户订单号
	ChannelNo     string `json:"channel_no"`     // 通道号
	OrderType     string `json:"order_type"`     // 订单类型：payin/payout
	MerchantNo    string `json:"merchant_no"`    // 商户号
	RetryCount    int    `json:"retry_count"`    // 重试次数
	CreatedAt     int64  `json:"created_at"`     // 任务创建时间
}

type QueryService struct {
	logger *zap.Logger
	cache  redis.Repo
}

var queryService *QueryService
var once sync.Once

// GetQueryService 获取查询服务实例（单例模式）
func GetQueryService(cache redis.Repo, logger *zap.Logger) *QueryService {
	once.Do(func() {
		queryService = &QueryService{
			logger: logger,
			cache:  cache,
		}
	})
	return queryService
}

func (s *QueryService) EnqueueQueryTask(task *OrderQueryTask) error {
	// 序列化任务
	taskData, err := json.Marshal(task)
	if err != nil {
		s.logger.Error("序列化通知任务失败",
			zap.Error(err),
			zap.String("order_no", task.OrderNo),
			zap.String("merchant_no", task.MerchantNo))
		return err
	}

	// 推送到Redis队列
	err = s.cache.Enqueue(OrderQueryQueue, string(taskData))
	if err != nil {
		s.logger.Error("推送查询任务到队列失败",
			zap.Error(err),
			zap.String("order_no", task.OrderNo),
			zap.String("merchant_no", task.MerchantNo))
		return err
	}

	return nil
}
