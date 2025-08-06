package authorized

import (
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/repository/mysql"
	"github.com/heigelove/cpay-payment/internal/repository/mysql/authorized"
)

func (s *service) Detail(ctx core.Context, id int32) (info *authorized.Authorized, err error) {
	qb := authorized.NewQueryBuilder()
	qb.WhereIsDeleted(mysql.EqualPredicate, -1)
	qb.WhereId(mysql.EqualPredicate, id)

	info, err = qb.First(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}
