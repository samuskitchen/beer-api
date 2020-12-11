package v1

import (
	"beer-api/domain/beer/application"
	"beer-api/domain/beer/application/v1/response"
	"beer-api/domain/beer/domain/model"
	repoDomain "beer-api/domain/beer/domain/repository"
	"beer-api/domain/beer/infrastructure/persistence"
	"beer-api/infrastructure/database"
	"beer-api/infrastructure/middleware"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

type BeersRouter struct {
	Repo repoDomain.BeersRepository
}

func NewBeerHandler(db *database.Data) *BeersRouter  {
	return &BeersRouter{
		Repo: persistence.NewBeersRepository(db),
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

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("the id is not a numeric").Error())
		return
	}

	ctx := r.Context()
	beerResult, err := br.Repo.GetBeerById(ctx, id)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
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
	ctx := r.Context()

	idStr := chi.URLParam(r, "beerID")
	currencyStr := r.URL.Query().Get("currency")
	quantityStr := r.URL.Query().Get("quantity")

	if idStr == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get id").Error())
		return
	}

	if currencyStr == "" {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("cannot get currency").Error())
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("the id is not a numeric").Error())
		return
	}

	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, errors.New("the quantity is not a numeric").Error())
		return
	}

	if quantity == 0 {
		quantity = 6
	}


	beerResult, err := br.Repo.GetBeerById(ctx, id)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, err.Error())
		return
	}

	if (model.Beers{}) == beerResult {
		_ = middleware.HTTPError(w, r, http.StatusNotFound, errors.New("the beer ID does not exist").Error())
		return
	}

	valueCurrency, err := application.GetCurrency(currencyStr)
	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	total := valueCurrency * float64(quantity)

	totalResponse := response.PriceResponse{
		PriceTotal: total,
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
	now := time.Now()
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

	ctx := r.Context()
	beers.CreatedAt = now

	if beers.ID != 0 {
		err = br.Repo.CreateBeerWithId(ctx, &beers)
	}else {
		err = br.Repo.CreateBeerWithOutId(ctx, &beers)
	}

	if err != nil {
		_ = middleware.HTTPError(w, r, http.StatusConflict, err.Error())
		return
	}

	_ = middleware.JSONMessages(w, r, http.StatusCreated, "Beer created")
}