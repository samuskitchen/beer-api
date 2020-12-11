package persistence

import (
	"beer-api/domain/beer/domain/model"
	repoDomain "beer-api/domain/beer/domain/repository"
	"beer-api/infrastructure/database"
	"context"
	"database/sql"
	"errors"
)

type sqlBeersRepository struct {
	Conn *database.Data
}

func NewBeersRepository(Connection *database.Data) repoDomain.BeersRepository {
	return &sqlBeersRepository{
		Conn: Connection,
	}
}

func (sb *sqlBeersRepository) GetAllBeers(ctx context.Context) ([]model.Beers, error) {
	rows, err := sb.Conn.DB.QueryContext(ctx, selectAllBeers)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var beers []model.Beers
	for rows.Next() {
		var beerRow model.Beers
		_ = rows.Scan(&beerRow.ID, &beerRow.Name, &beerRow.Brewery, &beerRow.Country, &beerRow.Price, &beerRow.Currency, &beerRow.CreatedAt)
		beers = append(beers, beerRow)
	}

	return beers, nil
}

func (sb *sqlBeersRepository) GetBeerById(ctx context.Context, id uint64) (model.Beers, error) {
	row := sb.Conn.DB.QueryRowContext(ctx, selectBeerById, id)

	var beerScan model.Beers
	err := row.Scan(&beerScan.ID, &beerScan.Name, &beerScan.Brewery, &beerScan.Country, &beerScan.Price, &beerScan.Currency, &beerScan.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows{
			return model.Beers{}, nil
		}

		return model.Beers{}, err
	}

	return beerScan, nil
}

func (sb *sqlBeersRepository) CreateBeerWithId(ctx context.Context, beers *model.Beers) error {
	beer, err := sb.GetBeerById(ctx, beers.ID)
	if err != nil {
		return err
	}

	if (model.Beers{}) != beer {
		return errors.New("beer ID already exists")
	}

	stmt, err := sb.Conn.DB.PrepareContext(ctx, insertBeerWithId)
	if err != nil {
		return err
	}

	row := stmt.QueryRowContext(ctx, &beers.ID, &beers.Name, &beers.Brewery, &beers.Country, &beers.Price, &beers.Currency, &beers.CreatedAt)
	defer stmt.Close()

	err = row.Scan(&beers.ID)
	if err != nil {
		return err
	}

	return nil
}

func (sb *sqlBeersRepository) CreateBeerWithOutId(ctx context.Context, beers *model.Beers) error {
	stmt, err := sb.Conn.DB.PrepareContext(ctx, insertBeerWithOutId)
	if err != nil {
		return err
	}

	row := stmt.QueryRowContext(ctx, &beers.Name, &beers.Brewery, &beers.Country, &beers.Price, &beers.Currency, &beers.CreatedAt)

	defer stmt.Close()

	err = row.Scan(&beers.ID)
	if err != nil {
		return err
	}

	return nil
}
