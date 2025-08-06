package notify

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/heigelove/cpay-payment/internal/repository/redis"

	"go.uber.org/zap"
)

const (
	// 队列名称
	MerchantNotifyQueueName = "merchant_notify_queue"

	// 最大重试次数
	MaxRetryCount = 3

	// 重试间隔（1分钟）
	RetryInterval = 1 * time.Minute

	// 分片数量，用于将消息分散到多个队列，提高并发处理能力
	QueueShardCount = 10

	// 单个worker的批处理大小
	BatchSize = 20

	// HTTP请求超时时间
	HTTPTimeout = 5 * time.Second

	// 每个worker处理的最大并发请求数
	MaxConcurrentRequests = 50
)

// 全局通知服务实例
var notifyService *NotifyService
var once sync.Once

// MerchantNotifyTask 商户通知任务
type MerchantNotifyTask struct {
	OrderNo      string     `json:"order_no"`       // 订单号
	MerchantNo   string     `json:"merchant_no"`    // 商户号
	NotifyURL    string     `json:"notify_url"`     // 通知URL
	Params       url.Values `json:"params"`         // 通知参数
	RetryCount   int        `json:"retry_count"`    // 重试次数
	OrderType    string     `json:"order_type"`     // 订单类型：payin or payout
	NextNotifyAt time.Time  `json:"next_notify_at"` // 下次通知时间
}

type NotifyService struct {
	logger *zap.Logger
	cache  redis.Repo
}

// GetNotifyService 获取通知服务实例（单例模式）
func GetNotifyService(cache redis.Repo, logger *zap.Logger) *NotifyService {
	once.Do(func() {
		notifyService = &NotifyService{
			logger: logger,
			cache:  cache,
		}
	})
	return notifyService
}

func NewNotifyService(cache redis.Repo, logger *zap.Logger) *NotifyService {
	return &NotifyService{
		logger: logger,
		cache:  cache,
	}
}

func (s *NotifyService) EnqueueNotifyTask(task *MerchantNotifyTask) error {
	// 序列化任务
	taskData, err := json.Marshal(task)
	if err != nil {
		s.logger.Error("序列化通知任务失败",
			zap.Error(err),
			zap.String("order_no", task.OrderNo),
			zap.String("merchant_no", task.MerchantNo))
		return err
	}

	// 根据订单号选择队列分片
	queueName := s.getQueueName(task.OrderNo)

	// 推送到Redis队列
	err = s.cache.Enqueue(queueName, string(taskData))
	if err != nil {
		s.logger.Error("推送通知任务到队列失败",
			zap.Error(err),
			zap.String("order_no", task.OrderNo),
			zap.String("merchant_no", task.MerchantNo))
		return err
	}

	return nil
}

// getQueueName 根据订单号哈希计算分片队列名
func (s *NotifyService) getQueueName(orderNo string) string {
	// 使用订单号来计算一致性哈希，确定队列分片
	// 这里使用简单的字符串哈希来做分片
	sum := 0
	for _, c := range orderNo {
		sum += int(c)
	}
	shardIndex := sum % QueueShardCount
	return fmt.Sprintf("%s_%d", MerchantNotifyQueueName, shardIndex)
}
