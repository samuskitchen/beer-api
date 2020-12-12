package integration

import (
	"beer-api/domain/beer/application/v1/response"
	"beer-api/domain/beer/domain/model"
	"beer-api/infrastructure/middleware"
	"beer-api/test/integration/seed"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// dataBeers is data for test
func dataBeers() []model.Beers {
	now := time.Now()

	return []model.Beers{
		{
			ID:        uint(1),
			Name:      "Golden",
			Brewery:   "Kross",
			Country:   "Chile",
			Price:     10.5,
			Currency:  "CLP",
			CreatedAt: now,
		},
		{
			ID:        uint(2),
			Name:      "Club Colombia",
			Brewery:   "Bavaria",
			Country:   "Colombia",
			Price:     2550,
			Currency:  "COP",
			CreatedAt: now,
		},
	}
}

func TestIntegration_GetAllBeersHandler(t *testing.T) {

	t.Run("No Content (no seed data)", func(tt *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/v1/beers/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusNotFound, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var beersResponse middleware.ErrorMessage
		if err := json.Unmarshal([]byte(w.Body.String()), &beersResponse); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if beersResponse.Message != "beers not found" {
			tt.Errorf("expected message to be returned, got %v", beersResponse.Message)
		}
	})

	t.Run("Ok (database has been seeded)", func(tt *testing.T) {
		defer func() {
			if err := seed.Truncate(dataConnection.DB); err != nil {
				t.Errorf("error truncating test database tables: %v", err)
			}
		}()

		expectedBeers, err := seed.BeersSeed(dataConnection.DB)
		if err != nil {
			t.Fatalf("error seeding beers: %v", err)
		}

		req, err := http.NewRequest(http.MethodGet, "/api/v1/beers/", nil)
		if err != nil {
			tt.Errorf("error creating request: %v", err)
		}

		w := httptest.NewRecorder()
		server.ServeHTTP(w, req)

		if e, a := http.StatusOK, w.Code; e != a {
			tt.Errorf("expected status code: %v, got status code: %v", e, a)
		}

		var beers []model.Beers
		if err := json.NewDecoder(w.Body).Decode(&beers); err != nil {
			tt.Errorf("error decoding response body: %v", err)
		}

		if d := cmp.Diff(expectedBeers[0].ID, beers[0].ID); d != "" {
			tt.Errorf("unexpected difference in response body:\n%v", d)
		}
	})
}

func TestIntegration_GetOneHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedBeers, err := seed.BeersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding beers: %v", err)
	}

	tests := []struct {
		Name         string
		BeerID       uint
		ExpectedBody model.Beers
		ExpectedCode int
	}{
		{
			Name:         "Get One Beer Successful",
			BeerID:       expectedBeers[0].ID,
			ExpectedBody: expectedBeers[0],
			ExpectedCode: http.StatusOK,
		},
		{
			Name:         "Beer Not Found",
			BeerID:       0,
			ExpectedBody: model.Beers{},
			ExpectedCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/beers/%d", test.BeerID), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var beersResponse model.Beers

				if err := json.NewDecoder(w.Body).Decode(&beersResponse); err != nil {
					t.Errorf("error decoding beersResponse body: %v", err)
				}

				if e, a := test.ExpectedBody.ID, beersResponse.ID; e != a {
					t.Errorf("expected user ID: %v, got user ID: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}

func TestIntegration_GetOneBoxPriceHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	expectedBeers, err := seed.BeersSeed(dataConnection.DB)
	if err != nil {
		t.Fatalf("error seeding beers: %v", err)
	}

	tests := []struct {
		Name          string
		BeerID        uint
		Currency      string
		Quantity      int
		ExpectedValue float64
		ExpectedCode  int
	}{
		{
			Name:          "Get One Beer Box Price Successful",
			BeerID:        expectedBeers[1].ID,
			Currency:      "USD",
			Quantity:      0,
			ExpectedValue: 4.47,
			ExpectedCode:  http.StatusOK,
		},
		{
			Name:          "Beer Not Found",
			Currency:      "USD",
			BeerID:        0,
			ExpectedValue: 0,
			ExpectedCode:  http.StatusNotFound,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/beers/%d/boxprice?currency=%s&quantity=%d", test.BeerID, test.Currency, test.Quantity), nil)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusNotFound {
				var priceResponse response.PriceResponse

				if err := json.NewDecoder(w.Body).Decode(&priceResponse); err != nil {
					t.Errorf("error decoding priceResponse body: %v", err)
				}

				if e, a := test.ExpectedValue, priceResponse.PriceTotal; e != a {
					t.Errorf("expected total value: %v, got total value: %v", e, a)
				}
			}
		}

		t.Run(test.Name, fn)
	}

}

func TestIntegration_CreateHandler(t *testing.T) {

	defer func() {
		if err := seed.Truncate(dataConnection.DB); err != nil {
			t.Errorf("error truncating test database tables: %v", err)
		}
	}()

	dataTest := dataBeers()[1]
	dataTest.ID = 0

	tests := []struct {
		Name         string
		RequestBody  model.Beers
		ExpectedCode int
	}{
		{
			Name:         "Create Beer With Id Successful",
			RequestBody:  dataBeers()[0],
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:         "Break Unique Id Constraint",
			RequestBody:  dataBeers()[0],
			ExpectedCode: http.StatusConflict,
		},
		{
			Name:         "No Data User",
			RequestBody:  model.Beers{},
			ExpectedCode: http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tests {
		fn := func(t *testing.T) {
			var b bytes.Buffer
			if err := json.NewEncoder(&b).Encode(test.RequestBody); err != nil {
				t.Errorf("error encoding request body: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, "/api/v1/beers/", &b)
			if err != nil {
				t.Errorf("error creating request: %v", err)
			}

			defer func() {
				if err := req.Body.Close(); err != nil {
					t.Errorf("error encountered closing request body: %v", err)
				}
			}()

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			if e, a := test.ExpectedCode, w.Code; e != a {
				t.Errorf("expected status code: %v, got status code: %v", e, a)
			}

			if test.ExpectedCode != http.StatusConflict && test.ExpectedCode != http.StatusUnprocessableEntity{
				var successfullyMessage middleware.SuccessfullyMessage

				if err := json.NewDecoder(w.Body).Decode(&successfullyMessage); err != nil {
					t.Errorf("error decoding successfullyMessage body: %v", err)
				}

				if successfullyMessage.Message != "Beer created" {
					t.Errorf("expected message to be returned, got %v", successfullyMessage.Message)
				}
			}
		}

		t.Run(test.Name, fn)
	}
}
