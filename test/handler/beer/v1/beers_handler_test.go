package v1

import (
	"beer-api/domain/beer/domain/repository"
	"net/http"
	"reflect"
	"testing"
)

func TestBeersRouter_CreateHandler(t *testing.T) {
	type fields struct {
		Repo repository.BeersRepository
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BeersRouter{
				Repo: tt.fields.Repo,
			}
		})
	}
}

func TestBeersRouter_GetAllBeersHandler(t *testing.T) {
	type fields struct {
		Repo repository.BeersRepository
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BeersRouter{
				Repo: tt.fields.Repo,
			}
		})
	}
}

func TestBeersRouter_GetOneBoxPriceHandler(t *testing.T) {
	type fields struct {
		Repo repository.BeersRepository
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BeersRouter{
				Repo: tt.fields.Repo,
			}
		})
	}
}

func TestBeersRouter_GetOneHandler(t *testing.T) {
	type fields struct {
		Repo repository.BeersRepository
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			br := &BeersRouter{
				Repo: tt.fields.Repo,
			}
		})
	}
}

func TestNewBeerHandler(t *testing.T) {
	type args struct {
		db *database.Data
	}
	tests := []struct {
		name string
		args args
		want *BeersRouter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBeerHandler(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBeerHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}
