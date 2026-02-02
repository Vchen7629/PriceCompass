package handler

import (
	"backend/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v4/pgxpool"
)

type API struct {
	Pool 		*pgxpool.Pool
	Validate 	*validator.Validate
	PlatformClients	*service.PlatformClients
}

// handler for validating api input params
func NewAPI(pool *pgxpool.Pool, clients *service.PlatformClients) *API {
	return &API{
		Pool: 				pool,
		Validate: 			validator.New(),
		PlatformClients: 	clients,
	}
}
