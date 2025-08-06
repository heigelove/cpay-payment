package sign

import (
	"net/url"
	"time"

	"github.com/heigelove/cpay-payment/pkg/errors"
)

func (s *signature) Verify(ts int64, sign string, params url.Values) (ok bool, err error) {
	if time.Now().Unix()-ts > s.ttl {
		err = errors.Errorf("request expired")
		return
	}

	tmpSign, err := s.Generate(params)
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
