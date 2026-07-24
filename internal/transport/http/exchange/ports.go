package exchange

type ExchangeApp interface {
	StartSession() (string, error)
}
