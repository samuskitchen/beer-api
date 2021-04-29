package v1

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"beer-api/domain/beer/application/v1/response"
	"beer-api/domain/beer/domain/model"
	repoDomain "beer-api/domain/beer/domain/repository"
	"beer-api/domain/beer/infrastructure/external"
	"beer-api/domain/beer/infrastructure/persistence"
	"beer-api/infrastructure/database"
	"beer-api/infrastructure/middleware"

	"github.com/go-chi/chi"
)

type BeersRouter struct {
	Repo   repoDomain.BeersRepository
	Client repoDomain.CurrencyInterface
}

func NewBeerHandler(db *database.Data, connectionHttp *http.Client) *BeersRouter {
	return &BeersRouter{
		Repo:   persistence.NewBeersRepository(db),
		Client: external.NewCurrencyRepository(connectionHttp),
	}
}

// swagger:route GET /beers beer getAllBeers
//
// GetAllBeersHandler.
// List all the beers found in the database
//
//     produces:
//      - application/json
//
//	   schemes: http, https
//
//     responses:
//        200: SwaggerAllBeersResponse
//		  404: SwaggerErrorMessage
//
// GetAllBeersHandler response all the beers.
func (br *BeersRouter) GetAllBeersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	beers, err := br.Repo.GetAllBeers(ctx)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	if beers == nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, errors.New("beers not found").Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, beers)
}

// swagger:route GET /beers/{beerID} beer idBeerPath
//
// GetOneHandler.
// Search for a beer by its ID
//
//     produces:
//      - application/json
//
//	   schemes: http, https
//
//     responses:
//        200: SwaggerBeersResponse
//		  404: SwaggerErrorMessage
//
// GetOneHandler response one beer by id.
func (br *BeersRouter) GetOneHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "beerID")
	if idStr == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("the id is not a numeric").Error())
		return
	}

	ctx := r.Context()
	beerResult, err := br.Repo.GetBeerById(ctx, uint(id))
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	if (model.Beers{}) == beerResult {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, errors.New("the beer ID does not exist").Error())
		return
	}

	_ = middleware.JSON(w, r, http.StatusOK, beerResult)
}

// swagger:route GET /beers/{beerID}/boxprice beer idBeerBoxPricePath
//
// GetOneBoxPriceHandler.
// Get the price of a case of beer by its ID
//
//     produces:
//      - application/json
//
//	   schemes: http, https
//
//     responses:
//        200: SwaggerPriceResponse
//		  404: SwaggerErrorMessage
//
// GetOneBoxPriceHandler get the price of a case of beer by its id
func (br *BeersRouter) GetOneBoxPriceHandler(w http.ResponseWriter, r *http.Request) {
	quantity := 6
	ctx := r.Context()

	idStr := chi.URLParam(r, "beerID")
	currencyStr := r.URL.Query().Get("currency")
	quantityStr := r.URL.Query().Get("quantity")

	if idStr == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	if currencyStr == "" || len(currencyStr) < 3 {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get currency").Error())
		return
	}

	if quantityStr != "" || len(quantityStr) > 0 {
		quantityValue, err := strconv.Atoi(quantityStr)
		if err != nil {
			_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("the quantity is not a numeric").Error())
			return
		}

		if quantityValue != 0 {
			quantity = quantityValue
		}
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("the id is not a numeric").Error())
		return
	}

	beerResult, err := br.Repo.GetBeerById(ctx, uint(id))
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if (model.Beers{}) == beerResult {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, errors.New("the beer ID does not exist").Error())
		return
	}

	valueCurrency, err := br.Client.GetCurrency(currencyStr, beerResult.Currency)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	valueTotalBeer := beerResult.Price * float64(quantity)
	total := valueCurrency[0] / valueCurrency[1] * valueTotalBeer
	totalFloat := math.Round(total*100) / 100

	totalResponse := response.PriceResponse{
		PriceTotal: totalFloat,
	}

	_ = middleware.JSON(w, r, http.StatusOK, totalResponse)
}

// swagger:route POST /beers beer beersRequest
//
// CreateHandler.
// Enter a new beer
//
//     consumes:
//     - application/json
//
//     produces:
//      - application/json
//
//	   schemes: http, https
//
//     responses:
//        201: SwaggerSuccessfullyMessage
//		  400: SwaggerErrorMessage
//		  409: SwaggerErrorMessage
//		  422: SwaggerErrorMessage
//
// CreateHandler Create a new beer.
func (br *BeersRouter) CreateHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var beers model.Beers

	err := json.NewDecoder(r.Body).Decode(&beers)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	defer r.Body.Close()
	userErrors := beers.Validate()
	if len(userErrors) > 0 {
		_ = middleware.HTTPErrors(w, r, http.StatusUnprocessableEntity, userErrors)
		return
	}

	beers.CreatedAt = time.Now()
	err = br.Repo.CreateBeerWithId(ctx, &beers)

	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	_ = middleware.JSONMessages(w, r, http.StatusCreated, "Beer created")
}
