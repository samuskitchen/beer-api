package response

// PriceResponse struct
type PriceResponse struct {
	PriceTotal float64 `json:"price_total,omitempty"`
}

// SwaggerPriceResponse it is the response of the price of beers
// swagger:response SwaggerPriceResponse
type SwaggerPriceResponse struct {
	//in: body
	Body PriceResponse
}
