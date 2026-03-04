package db

import (
	"backend/pkg/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

// helper method to validate if a value exists in the column in users table
func ValidateExistsUserTable(ctx context.Context, tx pgx.Tx, column, val string) error {
	var valueExists bool

	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM users WHERE %s = $1)`, column)

	err := tx.QueryRow(ctx, query, val).Scan(&valueExists)
	if err != nil {
		return err
	}

	if valueExists {
		return fmt.Errorf("%s already exists", column)
	}	

	return nil
}


// Use the provided sessionToken from frontend request and attempt to fetch the userId
// and username if the sessionToken is valid
func (r *Repository) ValidateSession(sessionToken string) (int, string, error) {
	ctx := context.Background()
	var userId int
	var username string

	err := db.WithTransaction(ctx, r.pool, func(tx pgx.Tx) error {
		query := `SELECT id, username FROM sessions WHERE token = $1`

		err := tx.QueryRow(ctx, query, sessionToken).Scan(&userId, &username)
		if err != nil {
			return fmt.Errorf("error fetching session for user: %w", err)
		}

		return nil
	})

	if err != nil {
		return 0, "", err
	}

	return userId, username, nil
}