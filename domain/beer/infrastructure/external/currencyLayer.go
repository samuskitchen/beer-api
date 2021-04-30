package external

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"beer-api/domain/beer/domain/model"
	repoDomain "beer-api/domain/beer/domain/repository"
)

type currencyClientData struct {
	Client *http.Client
}

// NewCurrencyRepository constructor
func NewCurrencyRepository(Connection *http.Client) repoDomain.CurrencyInterface {
	return &currencyClientData{
		Client: Connection,
	}
}

//GetCurrency Method that returns the currency of payment and beer as the base currency of conversion to the dollar,
// the order of the returned values are: currencyPay, currencyBeer
func (cl *currencyClientData) GetCurrency(currencyPay, currencyBeer string) ([]float64, error) {
	var valueEmpty []float64
	accessKey := os.Getenv("ACCESS_KEY_CURRENCY")

	responseCurrency, err := cl.Client.Get(fmt.Sprintf("http://apilayer.net/api/live?access_key=%s&currencies=%s,%s&source=USD&format=1", accessKey, currencyPay, currencyBeer))
	if err != nil {
		return valueEmpty, err
	}

	defer responseCurrency.Body.Close()
	if responseCurrency.StatusCode != 200 {
		return valueEmpty, fmt.Errorf("status code error: %d %s", responseCurrency.StatusCode, responseCurrency.Status)
	}

	responseData, err := ioutil.ReadAll(responseCurrency.Body)
	if err != nil {
		return valueEmpty, err
	}

	var currencyLayer model.CurrencyLayer
	err = json.Unmarshal(responseData, &currencyLayer)
	if err != nil {
		return valueEmpty, err
	}

	values := make([]float64, 0)

	valueCurrencyPay, ok := currencyLayer.Quotes["USD"+currencyPay].(float64)
	if !ok {
		return valueEmpty, errors.New("error get currency to pay")
	}
	values = append(values, valueCurrencyPay)

	valueCurrencyBeer, ok := currencyLayer.Quotes["USD"+currencyBeer].(float64)
	if !ok {
		return valueEmpty, errors.New("error get currency of the beer")
	}
	values = append(values, valueCurrencyBeer)

	return values, err

}
