// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// CurrencyInterface is an autogenerated mock type for the CurrencyInterface type
type CurrencyInterface struct {
	mock.Mock
}

// GetCurrency provides a mock function with given fields: currency
func (_m *CurrencyInterface) GetCurrency(currency string) (float64, error) {
	ret := _m.Called(currency)

	var r0 float64
	if rf, ok := ret.Get(0).(func(string) float64); ok {
		r0 = rf(currency)
	} else {
		r0 = ret.Get(0).(float64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(currency)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}