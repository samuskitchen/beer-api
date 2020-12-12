package repository

type CurrencyInterface interface {
	GetCurrency(currency, currencyBeer string) ([]float64, error)
}