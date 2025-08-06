package authorized

import (
	"github.com/heigelove/cpay-payment/configs"
	"github.com/heigelove/cpay-payment/internal/pkg/core"
	"github.com/heigelove/cpay-payment/internal/repository/mysql"
	"github.com/heigelove/cpay-payment/internal/repository/mysql/authorized_api"
	"github.com/heigelove/cpay-payment/internal/repository/redis"

	"gorm.io/gorm"
)

func (s *service) DeleteAPI(ctx core.Context, id int32) (err error) {
	// 先查询 id 是否存在
	authorizedApiInfo, err := authorized_api.NewQueryBuilder().
		WhereIsDeleted(mysql.EqualPredicate, -1).
		WhereId(mysql.EqualPredicate, id).
		First(s.db.GetDbR().WithContext(ctx.RequestContext()))

	if err == gorm.ErrRecordNotFound {
		return nil
	}

	data := map[string]interface{}{
		"is_deleted":   1,
		"updated_user": ctx.SessionUserInfo().UserName,
	}

	qb := authorized_api.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), data)
	if err != nil {
		return err
	}

	s.cache.Del(configs.RedisKeyPrefixSignature+authorizedApiInfo.BusinessKey, redis.WithTrace(ctx.Trace()))
	return
}
