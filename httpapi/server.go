package httpapi

import (
	"io"
	"os"
	"time"

	adauth "github.com/korylprince/go-ad-auth/v3"
	"github.com/korylprince/httputil/auth"
	"github.com/korylprince/httputil/auth/ad"
	"github.com/korylprince/httputil/session"
	"github.com/korylprince/httputil/session/memory"
)

//Server represents shared resources
type Server struct {
	auth         auth.Auth
	sessionStore session.Store
	groupRoles   RoleMap
	apiKeyRoles  map[string]string
	output       io.Writer
}

//NewServer returns a new server with the given resources
func NewServer(config *Config) *Server {
	authConfig := &adauth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.SecurityType(),
	}

	auth := ad.New(authConfig, nil, config.GroupRoleMap.Groups())

	sessionStore := memory.New(time.Minute * time.Duration(config.SessionExpiration))

	return &Server{
		auth:         auth,
		sessionStore: sessionStore,
		groupRoles:   config.GroupRoleMap,
		apiKeyRoles:  config.APIKeyRoleMap,
		output:       os.Stdout,
	}
}
