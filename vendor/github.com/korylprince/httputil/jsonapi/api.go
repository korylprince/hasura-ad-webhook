package jsonapi

import (
	"database/sql"
	"io"
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/korylprince/httputil/auth"
	"github.com/korylprince/httputil/session"
)

//APIRouter is an API Router
type APIRouter struct {
	mux          *mux.Router
	output       io.Writer
	auth         auth.Auth
	sessionStore session.Store
	hook         AuthHookFunc
}

//ServeHTTP implements the http.Handler interface
func (r *APIRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

//Handle registers a ReturnHandlerFunc with the given parameters
func (r *APIRouter) Handle(method, path string, handler ReturnHandlerFunc, auth bool) {
	action := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	if auth {
		handler = withAuth(r.sessionStore, handler)
	}

	r.mux.Methods(method).Path(path).Handler(
		withLogging(action, r.output,
			withJSONResponse(
				handler)))
}

//HandleTX registers a TXReturnHandlerFunc with the given parameters
func (r *APIRouter) HandleTX(method, path string, db *sql.DB, handler TXReturnHandlerFunc, auth bool) {
	action := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()

	rHandler := WithTX(db, handler)
	if auth {
		rHandler = withAuth(r.sessionStore, rHandler)
	}

	r.mux.Methods(method).Path(path).Handler(
		withLogging(action, r.output,
			withJSONResponse(
				rHandler)))
}

//New returns a new APIRouter
func New(output io.Writer, auth auth.Auth, store session.Store, hook AuthHookFunc) *APIRouter {
	r := &APIRouter{
		mux:          mux.NewRouter(),
		output:       output,
		auth:         auth,
		sessionStore: store,
		hook:         hook,
	}
	r.mux.NotFoundHandler = NotFoundJSONHandler
	r.Handle("POST", "/auth", r.authenticate, false)
	return r
}
