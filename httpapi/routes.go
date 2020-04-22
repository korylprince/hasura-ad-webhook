package httpapi

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/korylprince/httputil/jsonapi"
	"github.com/korylprince/httputil/session"
)

//API is the current API version
const API = "1.0"
const apiPath = "/api/" + API

//Router returns a new router
func (s *Server) Router() http.Handler {
	r := mux.NewRouter()

	apiRouter := jsonapi.New(s.output, s.auth, s.sessionStore, nil)
	r.PathPrefix(apiPath).Headers("X-Authorization-Type", "API-Key").Handler(http.StripPrefix(apiPath, apiRouter))
	apiRouter.Handle("GET", "/webhook", s.apiWebhook, false)

	var hook = func(sess session.Session) (bool, interface{}, error) {
		return true, map[string][]string{"roles": s.UserRoles(sess)}, nil
	}

	userRouter := jsonapi.New(s.output, s.auth, s.sessionStore, hook)
	r.PathPrefix(apiPath).Handler(http.StripPrefix(apiPath, userRouter))
	userRouter.Handle("GET", "/webhook", s.userWebhook, true)

	return handlers.CombinedLoggingHandler(s.output, r)
}
