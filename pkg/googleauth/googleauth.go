package googleauth

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image/png"
	"net/url"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
)

// GenerateSecret 生成谷歌验证器密钥
func GenerateSecret() (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ToopPay",
		AccountName: "admin",
	})
	if err != nil {
		return "", err
	}
	return key.Secret(), nil
}

// GetQRCodeURL 获取谷歌验证器二维码URL
func GetQRCodeURL(secret, username string) string {
	return fmt.Sprintf("otpauth://totp/ToopPay:%s?secret=%s&issuer=ToopPay",
		url.QueryEscape(username), secret)
}

// ValidateCode 验证谷歌验证码
func ValidateCode(secret, code string) bool {
	// 转换为标准格式
	secret = strings.TrimSpace(secret)
	secret = strings.ToUpper(secret)
	// 移除空格，便于用户输入
	code = strings.ReplaceAll(code, " ", "")

	// 添加容错，验证前后30秒的验证码也视为有效
	for i := -1; i <= 1; i++ {
		// 前后各偏移30秒
		t := time.Now().Add(time.Duration(i*30) * time.Second)
		if computeCode(secret, t) == code {
			return true
		}
	}
	return false
}

// computeCode 计算给定时间的验证码
func computeCode(secret string, t time.Time) string {
	// 将密钥转换为字节
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return ""
	}

	// 计算时间间隔
	counter := uint64(t.Unix() / 30)

	// 转换为字节
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)

	// 计算HMAC-SHA1
	h := hmac.New(sha1.New, key)
	h.Write(buf)
	sum := h.Sum(nil)

	// 截断
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	// 取模得到6位验证码
	mod := int32(value % 1000000)
	return fmt.Sprintf("%06d", mod)
}

// GenerateQRCodeBase64 生成二维码的Base64编码
func GenerateQRCodeBase64(secret, username string) (string, error) {
	// qrUrl := GetQRCodeURL(secret, username)
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "ToopPay",
		AccountName: username,
		Secret:      []byte(secret),
	})
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return "", err
	}
	// Use png.Encode to write the image to the buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	// 将图片编码为base64
	base64Encoding := "data:image/png;base64,"
	base64Str := base64Encoding + base64.StdEncoding.EncodeToString(buf.Bytes())

	return base64Str, nil
}
