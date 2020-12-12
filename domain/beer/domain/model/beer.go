package model

import (
	"time"
)

// Data of Beers
// swagger:model
type Beers struct {
	ID        uint	    `json:"id,omitempty"`
	// Required: true
	Name      string    `json:"name,omitempty"`
	// Required: true
	Brewery   string    `json:"brewery,omitempty"`
	// Required: true
	Country   string    `json:"country,omitempty"`
	// Required: true
	Price     float64   `json:"price,omitempty"`
	// Required: true
	Currency  string    `json:"currency,omitempty"`
	CreatedAt time.Time `json:"-"`
}

// swagger:parameters idBeerPath
type SwaggerBeerID struct {
	// in: path
	// Required: true
	BeerId string `json:"beerID"`
}

// swagger:parameters idBeerBoxPricePath
type SwaggerBeerBoxPrice struct {
	// in: path
	// Required: true
	BeerId string `json:"beerID"`

	// in: query
	// Required: true
	Currency string `json:"currency"`

	// in: query
	Quantity string `json:"quantity"`
}

// Information from Beers
// swagger:parameters beersRequest
type SwaggerBeersRequest struct {
	// in: body
	Body Beers
}

// Beers It is the response of the all beers information
// swagger:response SwaggerAllBeersResponse
type SwaggerAllBeersResponse struct {
	// in: body
	Body []struct {
		ID        uint	    `json:"id,omitempty"`
		Name      string    `json:"name,omitempty"`
		Brewery   string    `json:"brewery,omitempty"`
		Country   string    `json:"country,omitempty"`
		Price     float64   `json:"price,omitempty"`
		Currency  string    `json:"currency,omitempty"`
	}
}

// Beers It is the response of the beer information
// swagger:response SwaggerBeersResponse
type SwaggerBeersResponse struct {
	// in: body
	Body struct {
		ID        uint	    `json:"id,omitempty"`
		Name      string    `json:"name,omitempty"`
		Brewery   string    `json:"brewery,omitempty"`
		Country   string    `json:"country,omitempty"`
		Price     float64   `json:"price,omitempty"`
		Currency  string    `json:"currency,omitempty"`
	}
}

func (b *Beers) Validate() map[string]string {
	var errorMessages = make(map[string]string)

	if b.ID == 0 {
		errorMessages["id_required"] = "Id is required or id invalid"
	}

	if b.Name == "" {
		errorMessages["name_required"] = "names is required"
	}

	if b.Brewery == "" {
		errorMessages["brewery_required"] = "brewery is required"
	}

	if b.Country == "" {
		errorMessages["country_required"] = "country is required"
	}

	if b.Price == 0 {
		errorMessages["price_password"] = "price is required and different of zero"
	}

	if b.Currency == "" || len(b.Currency) < 3{
		errorMessages["currency_required"] = "currency is required and it has to be a valid currency"
	}

	return errorMessages
}
