//go:build unit

package db_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"backend/pkg/db"
)

// unit tests for HandleDatabaseErrors function
func TestHandleDatabaseErrors(t *testing.T) {
	t.Run("throws correct http error for db errors", func(t *testing.T) {
		testCases := []struct {
			errorCode 		string
			expectedStatus 	int
			expectedBody	string
		}{
			{"23505", http.StatusConflict, "Duplicate entry"},
			{"23503", http.StatusBadRequest, "Referenced record not found"},
			{"42P01", http.StatusInternalServerError, "Table not found"},
			{"99999", http.StatusInternalServerError, "Database error"},
		}
		
		for _, tc := range testCases {
			recorder := httptest.NewRecorder()
			pgErr := &pgconn.PgError{Code: tc.errorCode}

			result := db.HandleDatabaseErrors(recorder, pgErr)

			assert.True(t, result, "Should return true for PG error code %s", tc.errorCode)                                                
			assert.Equal(t, tc.expectedStatus, recorder.Code, "Wrong status for error code %s", tc.errorCode)                              
			assert.Contains(t, recorder.Body.String(), tc.expectedBody, "Wrong message for error code %s", tc.errorCode)  
		}
	})

	t.Run("returns false for non-PG errors", func(t *testing.T) {
		recorder := httptest.NewRecorder()                                                                                                 
		regularErr := errors.New("some regular error")                                                                                     
																																			
		result := db.HandleDatabaseErrors(recorder, regularErr)                                                                            
																																			
		assert.False(t, result)                                                                                                            
		assert.Equal(t, http.StatusOK, recorder.Code)  
	})
} 
