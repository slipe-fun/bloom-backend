package exchange

import "github.com/slipe-fun/skid-backend/internal/transport/ws/types"

type ExchangeHandler struct {
	exchangeApp ExchangeApp
	wsHub       *types.Hub
}

func NewExchangeHandler(exchangeApp ExchangeApp) *ExchangeHandler {
	return &ExchangeHandler{
		exchangeApp: exchangeApp,
	}
}
