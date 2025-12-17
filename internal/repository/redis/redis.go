package redis

import (
	"strings"
	"time"

	"github.com/heigelove/cpay-payment/configs"
	"github.com/heigelove/cpay-payment/pkg/errors"
	"github.com/heigelove/cpay-payment/pkg/timeutil"
	"github.com/heigelove/cpay-payment/pkg/trace"

	"github.com/go-redis/redis/v7"
)

type Option func(*option)

type Trace = trace.T

type option struct {
	Trace *trace.Trace
	Redis *trace.Redis
}

func newOption() *option {
	return &option{}
}

var _ Repo = (*cacheRepo)(nil)

type Repo interface {
	i()
	Set(key, value string, ttl time.Duration, options ...Option) error
	Get(key string, options ...Option) (string, error)
	// HGetAll hash表获取全部
	HGetAll(key string) (map[string]string, error)
	// HMset hash表设置字段
	HMSet(key string, value map[string]interface{}) error
	// HGet 获取hash表单个值
	HGet(key string, field string) (string, error)
	// HSet hash表设置字段
	HSet(key string, field string, value string) error
	// HDelete hash表删除字段
	HDelete(key string, field string) (int64, error)
	// TTL 获取key的剩余时间
	TTL(key string) (time.Duration, error)
	Expire(key string, ttl time.Duration) bool
	ExpireAt(key string, ttl time.Time) bool
	Del(key string, options ...Option) bool
	Exists(keys ...string) bool
	Incr(key string, options ...Option) int64
	Close() error
	Version() string
	Enqueue(queueName string, value string) error
	// Redis Set operations
	SAdd(key string, members ...interface{}) (int64, error)
	SRem(key string, members ...interface{}) (int64, error)
	SMembers(key string) ([]string, error)
	SIsMember(key string, member interface{}) (bool, error)
	SCard(key string) (int64, error)
	// Redis Sorted Set operations
	ZIncrBy(key string, increment float64, member string) (float64, error)
	ZScore(key string, member string) (float64, error)
	ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error)
}

type cacheRepo struct {
	client *redis.Client
}

func New() (Repo, error) {
	client, err := redisConnect()
	if err != nil {
		return nil, err
	}

	return &cacheRepo{
		client: client,
	}, nil
}

func (c *cacheRepo) i() {}

func redisConnect() (*redis.Client, error) {
	cfg := configs.Get().Redis
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Pass,
		DB:           cfg.Db,
		MaxRetries:   cfg.MaxRetries,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, errors.Wrap(err, "ping redis err")
	}

	return client, nil
}

// Set set some <key,value> into redis
func (c *cacheRepo) Set(key, value string, ttl time.Duration, options ...Option) error {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "set"
			opt.Redis.Key = key
			opt.Redis.Value = value
			opt.Redis.TTL = ttl.Minutes()
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if err := c.client.Set(key, value, ttl).Err(); err != nil {
		return errors.Wrapf(err, "redis set key: %s err", key)
	}

	return nil
}

// Get get some key from redis
func (c *cacheRepo) Get(key string, options ...Option) (string, error) {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "get"
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	value, err := c.client.Get(key).Result()
	if err != nil {
		return "", errors.Wrapf(err, "redis get key: %s err", key)
	}

	return value, nil
}

// TTL get some key from redis
func (c *cacheRepo) TTL(key string) (time.Duration, error) {
	ttl, err := c.client.TTL(key).Result()
	if err != nil {
		return -1, errors.Wrapf(err, "redis get key: %s err", key)
	}

	return ttl, nil
}

// Expire expire some key
func (c *cacheRepo) Expire(key string, ttl time.Duration) bool {
	ok, _ := c.client.Expire(key, ttl).Result()
	return ok
}

// ExpireAt expire some key at some time
func (c *cacheRepo) ExpireAt(key string, ttl time.Time) bool {
	ok, _ := c.client.ExpireAt(key, ttl).Result()
	return ok
}

func (c *cacheRepo) Exists(keys ...string) bool {
	if len(keys) == 0 {
		return true
	}
	value, _ := c.client.Exists(keys...).Result()
	return value > 0
}

func (c *cacheRepo) Del(key string, options ...Option) bool {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "del"
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}

	if key == "" {
		return true
	}

	value, _ := c.client.Del(key).Result()
	return value > 0
}

func (c *cacheRepo) Incr(key string, options ...Option) int64 {
	ts := time.Now()
	opt := newOption()
	defer func() {
		if opt.Trace != nil {
			opt.Redis.Timestamp = timeutil.CSTLayoutString()
			opt.Redis.Handle = "incr"
			opt.Redis.Key = key
			opt.Redis.CostSeconds = time.Since(ts).Seconds()
			opt.Trace.AppendRedis(opt.Redis)
		}
	}()

	for _, f := range options {
		f(opt)
	}
	value, _ := c.client.Incr(key).Result()
	return value
}

// HGetAll hash表获取全部
func (c *cacheRepo) HGetAll(key string) (map[string]string, error) {
	return c.client.HGetAll(key).Result()
}

// HMSet hash表设置字段
func (c *cacheRepo) HMSet(key string, value map[string]interface{}) error {
	_, err := c.client.HMSet(key, value).Result()
	return err
}

// HGet 获取hash表单个值
func (c *cacheRepo) HGet(key string, field string) (string, error) {
	return c.client.HGet(key, field).Result()
}

// HSet hash表设置字段
func (c *cacheRepo) HSet(key string, field string, value string) error {
	_, err := c.client.HSet(key, field, value).Result()
	return err
}

func (c *cacheRepo) HSetWithTTL(key string, field string, value string, ttl time.Duration) error {
	if err := c.HSet(key, field, value); err != nil {
		return err
	}

	// 设置过期时间
	if ttl > 0 {
		if !c.Expire(key, ttl) {
			return errors.Errorf("set expire for key %s failed", key)
		}
	}

	return nil
}

// HDelete hash表删除字段
func (c *cacheRepo) HDelete(key string, field string) (int64, error) {
	return c.client.HDel(key, field).Result()
}

// Enqueue 向队列尾部添加元素
func (c *cacheRepo) Enqueue(queueName string, value string) error {
	return c.client.RPush(queueName, value).Err()
}

// Close close redis client
func (c *cacheRepo) Close() error {
	return c.client.Close()
}

// WithTrace 设置trace信息
func WithTrace(t Trace) Option {
	return func(opt *option) {
		if t != nil {
			opt.Trace = t.(*trace.Trace)
			opt.Redis = new(trace.Redis)
		}
	}
}

// Version redis server version
func (c *cacheRepo) Version() string {
	server := c.client.Info("server").Val()
	spl1 := strings.Split(server, "# Server")
	spl2 := strings.Split(spl1[1], "redis_version:")
	spl3 := strings.Split(spl2[1], "redis_git_sha1:")
	return spl3[0]
}

// SAdd 向Redis Set中添加一个或多个成员
func (c *cacheRepo) SAdd(key string, members ...interface{}) (int64, error) {
	result, err := c.client.SAdd(key, members...).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis sadd key: %s err", key)
	}
	return result, nil
}

// SRem 从Redis Set中移除一个或多个成员
func (c *cacheRepo) SRem(key string, members ...interface{}) (int64, error) {
	result, err := c.client.SRem(key, members...).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis srem key: %s err", key)
	}
	return result, nil
}

// SMembers 获取Redis Set中的所有成员
func (c *cacheRepo) SMembers(key string) ([]string, error) {
	result, err := c.client.SMembers(key).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "redis smembers key: %s err", key)
	}
	return result, nil
}

// SIsMember 判断成员是否在Redis Set中
func (c *cacheRepo) SIsMember(key string, member interface{}) (bool, error) {
	result, err := c.client.SIsMember(key, member).Result()
	if err != nil {
		return false, errors.Wrapf(err, "redis sismember key: %s err", key)
	}
	return result, nil
}

// SCard 获取Redis Set中成员的数量
func (c *cacheRepo) SCard(key string) (int64, error) {
	result, err := c.client.SCard(key).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis scard key: %s err", key)
	}
	return result, nil
}

// ZIncrBy 为有序集合中的成员增加分数
func (c *cacheRepo) ZIncrBy(key string, increment float64, member string) (float64, error) {
	result, err := c.client.ZIncrBy(key, increment, member).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis zincrby key: %s member: %s err", key, member)
	}
	return result, nil
}

// ZScore 获取有序集合中成员的分数
func (c *cacheRepo) ZScore(key string, member string) (float64, error) {
	result, err := c.client.ZScore(key, member).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis zscore key: %s member: %s err", key, member)
	}
	return result, nil
}

// ZRevRangeWithScores 按分数从高到低获取有序集合中的成员和分数
func (c *cacheRepo) ZRevRangeWithScores(key string, start, stop int64) ([]redis.Z, error) {
	result, err := c.client.ZRevRangeWithScores(key, start, stop).Result()
	if err != nil {
		return nil, errors.Wrapf(err, "redis zrevrange key: %s err", key)
	}
	return result, nil
}
