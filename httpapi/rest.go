package httpapi

import (
	"net/http"
	"strings"

	"github.com/korylprince/httputil/jsonapi"
)

func (s *Server) apiWebhook(r *http.Request) (int, interface{}) {
	type response struct {
		Role string `json:"X-Hasura-Role"`
	}

	header := strings.Split(r.Header.Get("Authorization"), " ")

	if len(header) != 2 || header[0] != "Bearer" {
		return http.StatusUnauthorized, nil
	}

	if role, ok := s.apiKeyRoles[header[1]]; ok {
		return http.StatusOK, &response{Role: role}
	}

	return http.StatusUnauthorized, nil
}

func (s *Server) userWebhook(r *http.Request) (int, interface{}) {
	type response struct {
		UserID string `json:"X-Hasura-User-Id"`
		Role   string `json:"X-Hasura-Role"`
	}

	sess := jsonapi.GetSession(r)
	role := r.Header.Get("X-Hasura-Role")
	for _, r := range s.UserRoles(sess) {
		if r == role {
			return http.StatusOK, &response{UserID: sess.Username(), Role: role}
		}
	}

	return http.StatusUnauthorized, nil
}
