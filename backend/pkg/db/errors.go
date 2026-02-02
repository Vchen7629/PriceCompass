package db

import (
	"errors"
	"net/http"

	"github.com/jackc/pgconn"
)

// Helper function to check if types of database errors and return a formatted http error
func HandleDatabaseErrors(w http.ResponseWriter, db_err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(db_err, &pgErr) {
		switch pgErr.Code {
		case "23505": // duplicate key error
			http.Error(w, "Duplicate entry", http.StatusConflict)
		case "23503": // get request cant find it
			http.Error(w, "Referenced record not found", http.StatusBadRequest)
		case "42P01":
			http.Error(w, "Table not found", http.StatusInternalServerError)
		default:
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return true
	}
	return false
}