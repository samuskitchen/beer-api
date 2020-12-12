package v1

import (
	"beer-api/domain/beer/domain/model"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "beer-api/domain/beer/application/v1"
	repoMock "beer-api/domain/beer/domain/repository/mocks"
)

// dataBeers is data for test
func dataBeers() []model.Beers {
	now := time.Now()

	return []model.Beers{
		{
			ID:        uint64(1),
			Name:      "Golden",
			Brewery:   "Kross",
			Country:   "Chile",
			Price:     10.5,
			Currency:  "CLP",
			CreatedAt: now,
		},
		{
			ID:        uint64(2),
			Name:      "Club Colomhia",
			Brewery:   "Bavaria",
			Country:   "Colombia",
			Price:     2550,
			Currency:  "COP",
			CreatedAt: now,
		},
	}
}

func TestBeersRouter_GetAllBeersHandler(t *testing.T) {

	t.Run("Error Get All Beers Handler", func(tt *testing.T) {

		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/", nil)
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("GetAllBeers", mock.Anything).Return(nil, errors.New("error trace test"))

		testBeersHandler.GetAllBeersHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Not Found) Get All Beers Handler", func(tt *testing.T) {

		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/", nil)
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("GetAllBeers", mock.Anything).Return(nil, nil)

		testBeersHandler.GetAllBeersHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get All Beers Handler", func(tt *testing.T) {

		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/", nil)
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("GetAllBeers", mock.Anything).Return(dataBeers(), nil)

		testBeersHandler.GetAllBeersHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})
}

func TestBeersRouter_GetOneHandler(t *testing.T) {

	t.Run("Error Param Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}", nil)

		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.GetOneHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(model.Beers{}, errors.New("error sql")).Once()

		testBeersHandler.GetOneHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error Parse Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "no")

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.GetOneHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Not Found) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(model.Beers{}, sql.ErrNoRows).Once()

		testBeersHandler.GetOneHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(dataBeers()[0], nil).Once()

		testBeersHandler.GetOneHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

}

func TestBeersRouter_GetOneBoxPriceHandler(t *testing.T) {

	t.Run("(Error Param beerID) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Error Param Currency) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Error Param Quantity) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		queryParam := newRequest.URL.Query()
		queryParam.Add("currency", "COP")
		queryParam.Add("quantity", "no")
		newRequest.URL.RawQuery = queryParam.Encode()

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Error Parse beerID) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "no")

		queryParam := newRequest.URL.Query()
		queryParam.Add("currency", "COP")
		queryParam.Add("quantity", "1")
		newRequest.URL.RawQuery = queryParam.Encode()

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Error Get Beer) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		queryParam := newRequest.URL.Query()
		queryParam.Add("currency", "COP")
		newRequest.URL.RawQuery = queryParam.Encode()

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(model.Beers{}, errors.New("error sql")).Once()

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(No Found Beer) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		queryParam := newRequest.URL.Query()
		queryParam.Add("currency", "COP")
		newRequest.URL.RawQuery = queryParam.Encode()

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(model.Beers{}, nil).Once()

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("(Error Get Currency) Get One Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		queryParam := newRequest.URL.Query()
		queryParam.Add("currency", "COP")
		newRequest.URL.RawQuery = queryParam.Encode()

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		mockCurrencyRepository := &repoMock.CurrencyInterface{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository, Client: mockCurrencyRepository}

		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(dataBeers()[0], nil).Once()
		mockCurrencyRepository.On("GetCurrency", mock.Anything).Return(float64(0), errors.New("error get currency")).Once()

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Get One Box Price Handler", func(tt *testing.T) {

		newRecorder := httptest.NewRecorder()
		newRequest := httptest.NewRequest(http.MethodGet, "/api/v1/beers/{beerID}/boxprice", nil)

		newRequestCtx := chi.NewRouteContext()
		newRequestCtx.URLParams.Add("beerID", "1")

		queryParam := newRequest.URL.Query()
		queryParam.Add("currency", "COP")
		newRequest.URL.RawQuery = queryParam.Encode()

		newRequest = newRequest.WithContext(context.WithValue(newRequest.Context(), chi.RouteCtxKey, newRequestCtx))
		mockRepository := &repoMock.BeersRepository{}
		mockCurrencyRepository := &repoMock.CurrencyInterface{}
		testBeersHandler := &v1.BeersRouter{Repo: mockRepository, Client: mockCurrencyRepository}

		mockRepository.On("GetBeerById", mock.Anything, mock.Anything).Return(dataBeers()[0], nil).Once()
		mockCurrencyRepository.On("GetCurrency", mock.Anything).Return(3420.45, nil).Once()

		testBeersHandler.GetOneBoxPriceHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

}

func TestBeersRouter_CreateHandler(t *testing.T) {

	t.Run("Error Body Create Handler", func(tt *testing.T) {

		newRequest := httptest.NewRequest(http.MethodPost, "/api/v1/beers/", bytes.NewReader(nil))
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.CreateHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Validate Create Handler", func(tt *testing.T) {

		var userTest = dataBeers()[0]
		userTest.Name = ""

		marshal, err := json.Marshal(userTest)
		assert.NoError(tt, err)

		newRequest := httptest.NewRequest(http.MethodPost, "/api/v1/beers/", bytes.NewReader(marshal))
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}

		testBeersHandler.CreateHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)

	})

	t.Run("Error SQL With ID Create Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataBeers()[0])
		assert.NoError(tt, err)

		newRequest := httptest.NewRequest(http.MethodPost, "/api/v1/beers/", bytes.NewReader(marshal))
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("CreateBeerWithId", mock.Anything, mock.Anything).Return(errors.New("error sql"))

		testBeersHandler.CreateHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Error SQL WithOut ID Create Handler", func(tt *testing.T) {
		dataTest := dataBeers()[0]
		dataTest.ID = 0

		marshal, err := json.Marshal(dataTest)
		assert.NoError(tt, err)

		newRequest := httptest.NewRequest(http.MethodPost, "/api/v1/beers/", bytes.NewReader(marshal))
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("CreateBeerWithOutId", mock.Anything, mock.Anything).Return(errors.New("error sql"))

		testBeersHandler.CreateHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

	t.Run("Create With ID Handler", func(tt *testing.T) {

		marshal, err := json.Marshal(dataBeers()[0])
		assert.NoError(tt, err)

		newRequest := httptest.NewRequest(http.MethodPost, "/api/v1/beers/", bytes.NewReader(marshal))
		newRecorder := httptest.NewRecorder()
		mockRepository := &repoMock.BeersRepository{}

		testBeersHandler := &v1.BeersRouter{Repo: mockRepository}
		mockRepository.On("CreateBeerWithId", mock.Anything, mock.Anything).Return(nil)

		testBeersHandler.CreateHandler(newRecorder, newRequest)
		mockRepository.AssertExpectations(tt)
	})

}
