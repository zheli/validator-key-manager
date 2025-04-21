package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Create a mock database with ping monitoring enabled
	mockDB, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer mockDB.Close()

	// Set up test cases
	tests := []struct {
		name           string
		mockPingError  error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "healthy database",
			mockPingError:  nil,
			expectedStatus: http.StatusOK,
			expectedBody:   "ok",
		},
		{
			name:           "unhealthy database",
			mockPingError:  sql.ErrConnDone,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "Database connection error: sql: connection is already closed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock expectations
			if tt.mockPingError != nil {
				mock.ExpectPing().WillReturnError(tt.mockPingError)
			} else {
				mock.ExpectPing()
			}

			// Create a test request
			req := httptest.NewRequest("GET", "/healthz", nil)
			w := httptest.NewRecorder()

			// Create a new router for the test
			r := chi.NewRouter()
			r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
				if err := mockDB.Ping(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Database connection error: " + err.Error()))
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			})

			// Serve the request
			r.ServeHTTP(w, req)

			// Check the response
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())

			// Verify that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
