package external

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"beer-api/domain/beer/domain/repository"
	"beer-api/domain/beer/infrastructure/external"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	currencyRepositoryMock repository.CurrencyInterface
)

func newMockCurrency() repository.CurrencyInterface {
	currencyRepositoryMock = external.NewCurrencyRepository(http.DefaultClient)
	return currencyRepositoryMock
}

func closeMockCurrency() {
	currencyRepositoryMock = nil
}

func Test_currencyClientData_GetCurrency(t *testing.T) {

	err := godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		fmt.Printf("Error getting env %v\n", err)
	}

	t.Run("(Error CurrencyPay) Get Currency Successful", func(tt *testing.T) {
		currencyMockInterface := newMockCurrency()
		defer func() {
			closeMockCurrency()
		}()

		result, err := currencyMockInterface.GetCurrency("", "COP")
		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("(Error CurrencyBeer) Get Currency Successful", func(tt *testing.T) {
		currencyMockInterface := newMockCurrency()
		defer func() {
			closeMockCurrency()
		}()

		result, err := currencyMockInterface.GetCurrency("EUR", "")
		assert.Error(tt, err)
		assert.Nil(tt, result)
	})

	t.Run("Get Currency Successful", func(tt *testing.T) {
		currencyMockInterface := newMockCurrency()
		defer func() {
			closeMockCurrency()
		}()

		result, err := currencyMockInterface.GetCurrency("EUR", "COP")
		assert.NoError(tt, err)
		assert.NotNil(tt, result)
	})

}
