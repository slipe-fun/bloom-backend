package exchange

import (
	"github.com/redis/go-redis/v9"
)

type ExchangeApp struct {
	sessionApp SessionApp
	users      UserRepo
	rdb        *redis.Client
}

func NewExchangeApp(
	sessionApp SessionApp,
	users UserRepo,
	rdb *redis.Client,
) *ExchangeApp {
	return &ExchangeApp{
		sessionApp: sessionApp,
		users:      users,
		rdb:        rdb,
	}
}
