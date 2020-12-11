package repository

type CurrencyInterface interface {
	GetCurrency(currency string) (float64, error)
}