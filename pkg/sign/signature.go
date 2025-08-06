package sign

import (
	"net/url"
)

var _ Signature = (*signature)(nil)

type Signature interface {
	i()

	// Generate 生成签名
	Generate(params url.Values) (sign string, err error)

	// Verify 验证签名
	Verify(ts int64, sign string, params url.Values) (ok bool, err error)
}

type signature struct {
	secret string
	ttl    int64
}

func New(secret string, ttl int64) Signature {
	return &signature{
		secret: secret,
		ttl:    ttl,
	}
}

func (s *signature) i() {}
