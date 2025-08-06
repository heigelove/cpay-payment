package orderutils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/heigelove/cpay-payment/pkg/errors"
)

var (
	// 创建一个本地随机数生成器
	rng           *rand.Rand
	lastTimestamp int64
	sequence      int64
	mutex         sync.Mutex
)

// GenerateOrderNo 生成唯一订单号
// prefix: 订单号前缀，如 P 代表支付订单，W 代表提现订单
// 格式: 前缀 + 年月日时分秒 + 3位随机数 + 3位序列号
// 例如: P2025062015350112001 (P + 20250620153501 + 120 + 001)
func GenerateOrderNo(prefix string) string {
	mutex.Lock()
	defer mutex.Unlock()

	// 获取当前时间
	now := time.Now()
	timestamp := now.Format("20060102150405") // 年月日时分秒
	currentTimestamp := now.UnixNano() / 1e6  // 转换为毫秒级时间戳

	// 如果是同一毫秒内，序列号递增
	if currentTimestamp == lastTimestamp {
		sequence++
	} else {
		// 不是同一毫秒内，重置序列号
		sequence = 0
		lastTimestamp = currentTimestamp
	}

	// 生成3位随机数
	randomNum := rng.Intn(1000)

	// 格式化为完整订单号: 前缀 + 时间戳 + 3位随机数 + 3位序列号
	return fmt.Sprintf("%s%s%03d%03d", prefix, timestamp, randomNum, sequence)
}

// GeneratePayInOrderNo 生成代收订单号
func GeneratePayInOrderNo() string {
	return GenerateOrderNo("P")
}

// GeneratePayOutOrderNo 生成代付订单号
func GeneratePayOutOrderNo() string {
	return GenerateOrderNo("W")
}

// init 初始化随机数种子
func init() {
	source := rand.NewSource(time.Now().UnixNano())
	rng = rand.New(source)
}

// 执行实际的通知动作，使用 POST form 格式发送请求
func DoNotify(notifyUrl string, formValues url.Values) (bool, error) {
	if notifyUrl == "" {
		return false, fmt.Errorf("通知地址为空")
	}

	// 发送POST请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.PostForm(notifyUrl, formValues)
	if err != nil {
		return false, fmt.Errorf("发送通知请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 一般认为返回 success 或 OK 就是通知成功
	buf := make([]byte, 256)
	n, _ := resp.Body.Read(buf)
	respBody := string(buf[:n])

	// 这里简单验证，实际业务可能需要根据不同商户的规则来判断
	respBody = strings.ToLower(respBody)
	return strings.Contains(respBody, "success") || strings.Contains(respBody, "ok"), nil
}

// 执行实际的通知动作，使用 POST json 格式发送请求
func DoNotifyJson(notifyUrl string, data map[string]interface{}) (bool, error) {
	if notifyUrl == "" {
		return false, fmt.Errorf("通知地址为空")
	}

	// 构建通知数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("构建通知数据失败: %w", err)
	}

	// 发送POST请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", notifyUrl, strings.NewReader(string(jsonData)))
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// 一般认为返回 success 或 OK 就是通知成功
	buf := make([]byte, 256)
	n, _ := resp.Body.Read(buf)
	respBody := string(buf[:n])

	// 这里简单验证，实际业务可能需要根据不同商户的规则来判断
	respBody = strings.ToLower(respBody)
	return strings.Contains(respBody, "success") || strings.Contains(respBody, "ok"), nil
}

// 生成签名
func GenerateSignature(params url.Values, secret string) (sign string, err error) {
	var parts []string
	for key, values := range params {
		for _, value := range values {
			parts = append(parts, key+"="+value)
		}
	}
	// 按字母顺序排序参数
	sort.Strings(parts)
	// 使用 & 连接
	tmpStr := strings.Join(parts, "&") + "&key=" + secret
	fmt.Println("Signature String:", tmpStr)
	sign = strings.ToUpper(Md5(tmpStr))
	fmt.Println("Signature:", sign)
	return sign, nil
}

// VerifySignature 验证签名
func VerifySignature(params url.Values, ts int64, sign string, secret string) (ok bool, err error) {
	if time.Now().Unix()-ts > 5 {
		err = errors.Errorf("request expired")
		return
	}

	tmpSign, err := GenerateSignature(params, secret)
	if err != nil {
		err = errors.Errorf("generate sign error %v", err)
		return
	}

	if tmpSign != sign {
		err = errors.Errorf("sign not match %s != %s", tmpSign, sign)
		return
	}

	ok = true
	return
}

// MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
