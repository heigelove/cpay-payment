package authorized

import (
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/repository/mysql"
	"github.com/heigelove/cpay-payment/internal/repository/mysql/authorized_api"
)

type SearchAPIData struct {
	BusinessKey string `json:"business_key"` // 调用方key
}

func (s *service) ListAPI(ctx core.Context, searchAPIData *SearchAPIData) (listData []*authorized_api.AuthorizedApi, err error) {

	qb := authorized_api.NewQueryBuilder()
	qb = qb.WhereIsDeleted(mysql.EqualPredicate, -1)

	if searchAPIData.BusinessKey != "" {
		qb.WhereBusinessKey(mysql.EqualPredicate, searchAPIData.BusinessKey)
	}

	listData, err = qb.
		OrderById(false).
		QueryAll(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}
