package jsonapi

import (
	"database/sql"
	"net/http"
)

type contextKey int

const (
	contextKeyLogData contextKey = iota
	contextKeySession
)

//ReturnHandlerFunc returns an HTTP status code and body for the given request
type ReturnHandlerFunc func(*http.Request) (int, interface{})

//TXReturnHandlerFunc returns an HTTP status code and body for the given request and database handle
type TXReturnHandlerFunc func(*http.Request, *sql.Tx) (int, interface{})
