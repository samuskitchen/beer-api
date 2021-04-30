package repository

// CurrencyInterface interface
type CurrencyInterface interface {
	GetCurrency(currency, currencyBeer string) ([]float64, error)
}