package httpapi

import (
	"github.com/korylprince/httputil/auth/ad"
	"github.com/korylprince/httputil/session"
)

func (s *Server) UserRoles(sess session.Session) []string {
	rolesMap := make(map[string]struct{})
	user := sess.(*ad.User)

	for _, group := range user.Groups {
		if roles, ok := s.groupRoles[group]; ok {
			for _, r := range roles {
				rolesMap[r] = struct{}{}
			}
		}
	}

	roles := make([]string, 0, len(rolesMap))
	for r := range rolesMap {
		roles = append(roles, r)
	}

	return roles
}
