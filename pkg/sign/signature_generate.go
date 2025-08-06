package sign

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
)

// 生成签名
func (s *signature) Generate(params url.Values) (sign string, err error) {
	tmpStr := params.Encode()
	fmt.Println(tmpStr)
	sign = strings.ToUpper(Md5(tmpStr + "&key=" + s.secret))
	fmt.Println(sign)
	return
}

// MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
