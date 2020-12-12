// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	model "beer-api/domain/beer/domain/model"
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// BeersRepository is an autogenerated mock type for the BeersRepository type
type BeersRepository struct {
	mock.Mock
}

// CreateBeerWithId provides a mock function with given fields: ctx, beers
func (_m *BeersRepository) CreateBeerWithId(ctx context.Context, beers *model.Beers) error {
	ret := _m.Called(ctx, beers)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.Beers) error); ok {
		r0 = rf(ctx, beers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllBeers provides a mock function with given fields: ctx
func (_m *BeersRepository) GetAllBeers(ctx context.Context) ([]model.Beers, error) {
	ret := _m.Called(ctx)

	var r0 []model.Beers
	if rf, ok := ret.Get(0).(func(context.Context) []model.Beers); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Beers)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBeerById provides a mock function with given fields: ctx, id
func (_m *BeersRepository) GetBeerById(ctx context.Context, id uint) (model.Beers, error) {
	ret := _m.Called(ctx, id)

	var r0 model.Beers
	if rf, ok := ret.Get(0).(func(context.Context, uint) model.Beers); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(model.Beers)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
