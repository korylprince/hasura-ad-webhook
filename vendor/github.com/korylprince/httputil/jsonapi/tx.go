package jsonapi

import (
	"database/sql"
	"fmt"
	"net/http"
)

//WithTX returns a ReturnHandlerFunc for the given database and TXReturnHandlerFunc
func WithTX(db *sql.DB, next TXReturnHandlerFunc) ReturnHandlerFunc {
	return func(r *http.Request) (int, interface{}) {
		tx, err := db.Begin()
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Unable to start database transaction: %v", err)
		}

		status, body := next(r, tx)

		if status != http.StatusOK {
			if err = tx.Rollback(); err != nil {
				if pErr, ok := body.(error); ok {
					return http.StatusInternalServerError, fmt.Errorf("Unable to rollback database transaction: %v; Previous error: HTTP %d %s: %v", err, status, http.StatusText(status), pErr)
				}
				return http.StatusInternalServerError, fmt.Errorf("Unable to rollback database transaction: %v", err)
			}
			return status, body
		}

		if err = tx.Commit(); err != nil {
			if pErr, ok := body.(error); ok {
				return http.StatusInternalServerError, fmt.Errorf("Unable to commit database transaction: %v; Previous error: HTTP %d %s: %v", err, status, http.StatusText(status), pErr)
			}
			return http.StatusInternalServerError, fmt.Errorf("Unable to commit database transaction: %v", err)
		}

		return status, body
	}
}
