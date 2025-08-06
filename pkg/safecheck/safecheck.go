package safecheck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ReqBodyKey        = "req-body"
	ResBodyKey        = "res-body"
	TreePathDelimiter = "."
)

type SafeCheckConfig struct {
	AllowedPathPrefixes []string
	SkippedPathPrefixes []string
}

// 检查SQL注入
func checkSQLInjection(c *gin.Context) error {
	// 1. 检查URL查询参数
	query := c.Request.URL.Query()
	if err := checkQueryParams(query); err != nil {
		return err
	}
	// 将查询参数重新编码到URL中
	// 以确保它们在请求中被正确处理
	// c.Request.URL.RawQuery = query.Encode()

	// 2. 检查POST表单数据
	if c.ContentType() == "application/x-www-form-urlencoded" {
		if err := c.Request.ParseForm(); err == nil {
			form := c.Request.PostForm // 只获取POST表单数据，不包括URL查询参数
			if err := checkQueryParams(form); err != nil {
				return err
			}
			// 重新设置表单数据（不覆盖URL查询参数）
			// c.Request.PostForm = form
		}
	}

	// 3. 检查JSON请求体
	bodyData, exists := c.Get(ReqBodyKey)
	if exists {
		if bodyBytes, ok := bodyData.([]byte); ok && len(bodyBytes) > 0 {
			var jsonData interface{}
			if err := json.Unmarshal(bodyBytes, &jsonData); err == nil {
				if err := checkJSONData(jsonData); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 检查URL查询参数
func checkQueryParams(values url.Values) error {
	for _, vals := range values {
		for _, v := range vals {
			if isSQLInjection(v) {
				return errors.New("Bad Request")
			}
		}
	}
	return nil
}

// 检查JSON数据
func checkJSONData(data interface{}) error {
	switch v := data.(type) {
	case map[string]interface{}:
		for _, value := range v {
			if err := checkJSONData(value); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, item := range v {
			if err := checkJSONData(item); err != nil {
				return err
			}
		}
	case string:
		if isSQLInjection(v) {
			return errors.New("Bad Request")
		}
	}
	return nil
}

// SQL注入检测正则表达式
var (
	// SQL关键字正则
	sqlKeywordRegex = regexp.MustCompile(`(?i)\b(SELECT|INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|TRUNCATE|REPLACE|UNION|JOIN|WHERE|HAVING|GROUP\s+BY|ORDER\s+BY|OR\s+1=1|AND\s+1=1)\b`)

	// SQL注释正则
	sqlCommentRegex = regexp.MustCompile(`(?i)(--|#|\/*|;--)`)

	// SQL批处理正则
	sqlBatchRegex = regexp.MustCompile(`(?i)(;[\s\w]*SELECT|;[\s\w]*INSERT|;[\s\w]*UPDATE|;[\s\w]*DELETE|;[\s\w]*DROP)`)

	// SQL字符串攻击正则
	sqlStringRegex = regexp.MustCompile(`(?i)('\s*OR\s*'[\s\w]*'='[\s\w]*'|"\s*OR\s*"[\s\w]*"="[\s\w]*")`)
)

// 检测字符串是否包含SQL注入风险
func isSQLInjection(value string) bool {
	// 如果值很短或为空，跳过检查
	if len(value) < 3 {
		return false
	}

	// 检查是否包含SQL关键字和特殊字符组合
	if sqlKeywordRegex.MatchString(value) && (strings.Contains(value, "'") || strings.Contains(value, "\"") || strings.Contains(value, ";")) {
		return true
	}

	// 检查是否包含SQL注释
	// if sqlCommentRegex.MatchString(value) {
	// 	return true
	// }

	// 检查是否包含SQL批处理命令
	if sqlBatchRegex.MatchString(value) {
		return true
	}

	// 检查是否包含常见的SQL字符串攻击
	if sqlStringRegex.MatchString(value) {
		return true
	}

	return false
}

func SafeCheckWithConfig(config SafeCheckConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !AllowedPathPrefixes(c, config.AllowedPathPrefixes...) ||
			SkippedPathPrefixes(c, config.SkippedPathPrefixes...) ||
			c.Request.Body == nil {
			c.Next()
			return
		}

		// 记录请求体
		if c.Request.Body != nil {
			bodyData, err := c.GetRawData()
			if err != nil {
				logger.Error("获取请求体失败", zap.Error(err))
				c.AbortWithStatusJSON(400, gin.H{"error": "Bad Request"})
				return
			}
			c.Set(ReqBodyKey, bodyData)

			// 重建请求体，以便后续处理程序可以使用
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyData))
		}

		//sql 注入检查， 策略注册检查下放到策略服务
		// SQL 注入检查
		if err := checkSQLInjection(c); err != nil {
			logger.Error(fmt.Sprint("检测到SQL注入风险: %s", err.Error()))
			c.Abort()
			return
		}

		c.Next()
	}

}
