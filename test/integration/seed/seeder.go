package seed

import (
	"beer-api/infrastructure/database"
	"database/sql"
	"github.com/pkg/errors"
	"time"

	modelBeer "beer-api/domain/beer/domain/model"
	_ "github.com/lib/pq"
)

// Open returns a new database connection for the test database.
func Open() *database.Data {
	return database.NewTest()
}

// Truncate removes all seed data from the test database.
func Truncate(dbc *sql.DB) error {
	stmt := "TRUNCATE TABLE beers RESTART IDENTITY CASCADE;"

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

// BeersSeed handles seeding the beers table in the database for integration tests.
func BeersSeed(dbc *sql.DB) ([]modelBeer.Beers, error) {
	now := time.Now()

	beers := []modelBeer.Beers{
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

	for i := range beers {
		query := `INSERT INTO beers (id, "name", brewery, country, price, currency, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`

		stmt, err := dbc.Prepare(query)
		if err != nil {
			return nil, errors.Wrap(err, "prepare beer insertion")
		}

		row := stmt.QueryRow(&beers[i].ID, &beers[i].Name, &beers[i].Brewery, &beers[i].Country, &beers[i].Price, &beers[i].Currency, &beers[i].CreatedAt)

		if err = row.Scan(&beers[i].ID); err != nil {
			if err := stmt.Close(); err != nil {
				return nil, errors.Wrap(err, "close psql statement")
			}

			return nil, errors.Wrap(err, "capture beer id")
		}

		if err := stmt.Close(); err != nil {
			return nil, errors.Wrap(err, "close psql statement")
		}
	}

	return beers, nil
}
