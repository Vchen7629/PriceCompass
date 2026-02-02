package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
)

type API struct {
	Pool 		*pgxpool.Pool
	Validate 	*validator.Validate
}

// handler for validating api input params
func NewAPI(pool *pgxpool.Pool) *API {
	return &API{
		Pool: 		pool,
		Validate: 	validator.New(),
	}
}
