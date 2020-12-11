package repository

import (
	"beer-api/domain/beer/domain/model"
	"context"
)

type BeersRepository interface {
	GetAllBeers(ctx context.Context) ([]model.Beers, error)
	GetBeerById(ctx context.Context, id uint64) (model.Beers, error)
	CreateBeerWithId(ctx context.Context, beers *model.Beers) error
	CreateBeerWithOutId(ctx context.Context, beers *model.Beers) error
}