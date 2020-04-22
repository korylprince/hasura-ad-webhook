package jsonapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/korylprince/httputil/session"
)

//AuthHookFunc is used to further check authorization of a Session and return extra attributes
type AuthHookFunc func(session.Session) (status bool, attrs interface{}, err error)

//withAuth returns a ReturnHandlerFunc that checks the Session of the given request.
func withAuth(store session.Store, next ReturnHandlerFunc) ReturnHandlerFunc {
	return func(r *http.Request) (int, interface{}) {
		header := strings.Split(r.Header.Get("Authorization"), " ")

		if len(header) != 2 || header[0] != "Bearer" || len(header[1]) != 36 {
			return http.StatusBadRequest, errors.New("Invalid Authorization header")
		}

		session, err := store.Read(header[1])
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Unexpected error when checking session id %s: %v", header[1], err)
		}

		if session == nil {
			return http.StatusUnauthorized, fmt.Errorf("Session doesn't exist for id %s", header[1])
		}

		(r.Context().Value(contextKeyLogData)).(*logData).User = session.Username()

		ctx := context.WithValue(r.Context(), contextKeySession, session)

		status, body := next(r.WithContext(ctx))
		return status, body
	}
}

//GetSession returns the Session for the given request
func GetSession(r *http.Request) session.Session {
	return (r.Context().Value(contextKeySession)).(session.Session)
}

func (router *APIRouter) authenticate(r *http.Request) (int, interface{}) {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type response struct {
		Username    string      `json:"username"`
		DisplayName string      `json:"display_name"`
		SessionID   string      `json:"session_id"`
		Attrs       interface{} `json:"attrs"`
	}

	req := new(request)

	if err := ParseJSONBody(r, req); err != nil {
		return http.StatusBadRequest, err
	}

	session, err := router.auth.Authenticate(req.Username, req.Password)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to authenticate: %v", err)
	}

	if session == nil {
		return http.StatusUnauthorized, errors.New("Invalid username or password")
	}

	resp := &response{
		Username:    session.Username(),
		DisplayName: session.DisplayName(),
	}

	if router.hook != nil {
		status, attrs, err := router.hook(session)
		if err != nil {
			return http.StatusInternalServerError, fmt.Errorf("Error when running authentication hook: %v", err)
		}

		if !status {
			return http.StatusUnauthorized, nil
		}
		resp.Attrs = attrs
	}

	id, err := router.sessionStore.Create(session)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("Unable to create session: %v", err)
	}

	resp.SessionID = id

	return http.StatusOK, resp
}
