package external

import (
	"beer-api/domain/beer/domain/model"
	repoDomain "beer-api/domain/beer/domain/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type httpClientData struct {
	Client *http.Client
}

func NewCurrencyRepository(Connection *http.Client) repoDomain.CurrencyInterface {
	return &httpClientData{
		Client: Connection,
	}
}

func (cl *httpClientData) GetCurrency(currency string) (float64, error) {
	accessKey := os.Getenv("ACCESS_KEY_CURRENCY")

	responseCurrency, err := cl.Client.Get("http://apilayer.net/api/live?access_key=" + accessKey + "&currencies=" + currency + "&source=USD&format=1")
	if err != nil {
		return 0, err
	}

	defer responseCurrency.Body.Close()
	if responseCurrency.StatusCode != 200 {
		return 0, fmt.Errorf("status code error: %d %s", responseCurrency.StatusCode, responseCurrency.Status)
	}

	responseData, err := ioutil.ReadAll(responseCurrency.Body)
	if err != nil {
		return 0, err
	}

	var currencyLayer model.CurrencyLayer
	err = json.Unmarshal(responseData, &currencyLayer)
	if err != nil {
		return 0, err
	}

	value, ok := currencyLayer.Quotes["USD"+currency].(float64)
	if !ok {
		return 0, errors.New("error get currency")
	}

	return value, err

}
