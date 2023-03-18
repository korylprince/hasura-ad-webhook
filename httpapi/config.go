package httpapi

import (
	"fmt"
	"log"
	"strings"

	auth "github.com/korylprince/go-ad-auth/v3"
)

// RoleMap represents groups mapped to roles
type RoleMap map[string][]string

// Decode parses the format k1:v1,v2,...;k2:v3,v4,... to map[string][]string
func (r *RoleMap) Decode(val string) error {
	m := make(RoleMap)
	for _, kv := range strings.Split(val, ";") {
		sp := strings.Split(kv, ":")
		if len(sp) != 2 {
			return fmt.Errorf("Unable to parse group string: %s", kv)
		}
		m[sp[0]] = strings.Split(sp[1], ",")
	}

	*r = m
	return nil
}

// Groups returns the groups in r
func (r *RoleMap) Groups() []string {
	groups := make([]string, 0, len(*r))
	for g := range *r {
		groups = append(groups, g)
	}

	return groups
}

// Config represents options given in the environment
type Config struct {
	SessionExpiration int `default:"60"` //in minutes

	LDAPServer   string `required:"true"`
	LDAPPort     int    `default:"389" required:"true"`
	LDAPBaseDN   string `required:"true"`
	LDAPSecurity string `default:"none" required:"true"`

	GroupRoleMap  RoleMap           `required:"true"` //map[group]roles
	APIKeyRoleMap map[string]string //map[api_key]role

	ListenAddr string `default:":8080" required:"true"` //addr format used for net.Dial; required
	Prefix     string //url prefix to mount api to without trailing slash
}

// SecurityType returns the auth.SecurityType for the config
func (c *Config) SecurityType() auth.SecurityType {
	switch strings.ToLower(c.LDAPSecurity) {
	case "", "none":
		return auth.SecurityNone
	case "tls":
		return auth.SecurityTLS
	case "starttls":
		return auth.SecurityStartTLS
	default:
		log.Fatalln("Invalid LDAPSECURITY:", c.LDAPSecurity)
	}
	panic("unreachable")
}
