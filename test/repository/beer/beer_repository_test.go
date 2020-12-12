package beer

import (
	"beer-api/domain/beer/domain/model"
	"beer-api/domain/beer/domain/repository"
	"beer-api/domain/beer/infrastructure/persistence"
	"beer-api/infrastructure/database"
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

// represent the repository
var (
	dbMockBeers        *sql.DB
	connMockBeer       database.Data
	beerRepositoryMock repository.BeersRepository
)

// newMockBeer initialize mock connection to database
func newMockBeer() sqlmock.Sqlmock {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	dbMockBeers = db
	connMockBeer = database.Data{
		DB: dbMockBeers,
	}

	beerRepositoryMock = persistence.NewBeersRepository(&connMockBeer)

	return mock
}

// closeMockBeer Close attaches the provider and close the connection
func closeMockBeer() {
	err := dbMockBeers.Close()
	if err != nil {
		log.Println("Error close database test")
	}
}

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

func Test_sqlBeersRepository_GetAllBeers(t *testing.T) {

	t.Run("Error SQL", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		mock.ExpectQuery("SELECT 1 FROM beers")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		beers, err := beerRepositoryMock.GetAllBeers(ctx)
		assert.Error(tt, err)
		assert.Nil(tt, beers)
	})

	t.Run("Get All Beers Successful", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		beersData := dataBeers()
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beersData[0].ID, beersData[0].Name, beersData[0].Brewery, beersData[0].Country, beersData[0].Price, beersData[0].Currency, beersData[0].CreatedAt).
			AddRow(beersData[1].ID, beersData[1].Name, beersData[1].Brewery, beersData[1].Country, beersData[1].Price, beersData[1].Currency, beersData[1].CreatedAt)

		mock.ExpectQuery(selectAllBeersTest).WillReturnRows(rows)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		beers, err := beerRepositoryMock.GetAllBeers(ctx)
		assert.NotEmpty(tt, beers)
		assert.NoError(tt, err)
		assert.Len(tt, beers, 2)
	})
}

func Test_sqlBeersRepository_GetBeerById(t *testing.T) {
	beerTest := dataBeers()[0]

	t.Run("Error SQL", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(nil).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		beer, err := beerRepositoryMock.GetBeerById(ctx, 0)
		assert.Error(tt, err)
		assert.NotNil(tt, beer)
	})

	t.Run("Not Found", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt).
			RowError(0, sql.ErrNoRows)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		beer, err := beerRepositoryMock.GetBeerById(ctx, beerTest.ID)
		assert.NoError(tt, err)
		assert.NotNil(tt, beer)
		assert.Equal(tt, model.Beers{}, beer)
	})

	t.Run("Get Beer By Id Successful", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		beer, err := beerRepositoryMock.GetBeerById(ctx, beerTest.ID)
		assert.NoError(tt, err)
		assert.NotNil(tt, beer)
	})

}

func Test_sqlBeersRepository_CreateBeerWithId(t *testing.T) {
	beerTest := dataBeers()[0]
	beerEmptyTest := model.Beers{}
	beerErrorTest := dataBeers()[0]

	t.Run("Error Get Beer", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt).
			RowError(0, errors.New("error"))

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := beerRepositoryMock.CreateBeerWithId(ctx, &dataBeers()[0])
		assert.Error(tt, err)
	})

	t.Run("Beer Found", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerTest.ID).WillReturnRows(row)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := beerRepositoryMock.CreateBeerWithId(ctx, &dataBeers()[0])
		assert.Error(tt, err)
	})

	t.Run("Error SQL", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerEmptyTest.ID, beerEmptyTest.Name, beerEmptyTest.Brewery, beerEmptyTest.Country, beerEmptyTest.Price, beerEmptyTest.Currency, beerEmptyTest.CreatedAt)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerEmptyTest.ID).WillReturnRows(row)

		prep := mock.ExpectPrepare("insertBeerWithIdTest")
		prep.ExpectExec().
			WithArgs(beerErrorTest.ID, beerErrorTest.Name, beerErrorTest.Brewery, beerErrorTest.Country, beerErrorTest.Price, beerErrorTest.Currency, beerErrorTest.CreatedAt).
			WillReturnResult(sqlmock.NewResult(0, 0))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := beerRepositoryMock.CreateBeerWithId(ctx, &beerEmptyTest)
		assert.Error(tt, err)
	})

	t.Run("Error Scan Row", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerEmptyTest.ID, beerEmptyTest.Name, beerEmptyTest.Brewery, beerEmptyTest.Country, beerEmptyTest.Price, beerEmptyTest.Currency, beerEmptyTest.CreatedAt)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerEmptyTest.ID).WillReturnRows(row)

		prep := mock.ExpectPrepare(insertBeerWithIdTest)
		prep.ExpectQuery().
			WithArgs(beerTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"first_name"}).AddRow("Error"))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := beerRepositoryMock.CreateBeerWithId(ctx, &beerEmptyTest)
		assert.Error(tt, err)
	})

	t.Run("Create Beer With Id Successful", func(tt *testing.T) {
		mock := newMockBeer()
		defer func() {
			closeMockBeer()
		}()

		row := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at"}).
			AddRow(beerEmptyTest.ID, beerEmptyTest.Name, beerEmptyTest.Brewery, beerEmptyTest.Country, beerEmptyTest.Price, beerEmptyTest.Currency, beerEmptyTest.CreatedAt)

		mock.ExpectQuery(selectBeerByIdTest).WithArgs(beerEmptyTest.ID).WillReturnRows(row)

		prep := mock.ExpectPrepare(insertBeerWithIdTest)
		prep.ExpectQuery().
			WithArgs(beerEmptyTest.ID, beerTest.Name, beerTest.Brewery, beerTest.Country, beerTest.Price, beerTest.Currency, beerTest.CreatedAt).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(beerEmptyTest.ID))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		beerTest.ID = beerEmptyTest.ID
		err := beerRepositoryMock.CreateBeerWithId(ctx, &beerTest)
		assert.NoError(tt, err)
	})

}