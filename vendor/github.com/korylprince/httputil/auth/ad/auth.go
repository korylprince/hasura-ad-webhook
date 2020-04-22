package ad

import (
	"fmt"

	ldap "github.com/go-ldap/ldap/v3"
	adauth "github.com/korylprince/go-ad-auth/v3"
	"github.com/korylprince/httputil/session"
)

//User represents an Active Directory User
type User struct {
	username string
	Entry    *ldap.Entry
	Groups   []string
}

//Username returns the User's username
func (u *User) Username() string {
	return u.username
}

//DisplayName returns the User's display name
func (u *User) DisplayName() string {
	return u.Entry.GetAttributeValue("displayName")
}

//Auth represents an Active Directory authentication mechanism
type Auth struct {
	config *adauth.Config
	attrs  []string
	groups []string
}

//New returns a new *Auth with the given configuration, user attributes, and groups
func New(config *adauth.Config, attrs, groups []string) *Auth {
	for _, attr := range attrs {
		if attr == "displayName" {
			return &Auth{config: config, attrs: attrs, groups: groups}
		}
	}

	return &Auth{config: config, attrs: append(attrs, "displayName"), groups: groups}
}

//Authenticate authenticates the given credentials and returns the User associated with the account if successful,
//or nil if not. If an error occurs it is returned.
func (a *Auth) Authenticate(username, password string) (user session.Session, err error) {
	status, entry, groups, err := adauth.AuthenticateExtended(a.config, username, password, a.attrs, a.groups)
	if err != nil {
		return nil, fmt.Errorf("Error attempting to authenticate as %s: %v", username, err)
	}

	if !status {
		return nil, nil
	}

	if len(a.groups) > 0 && len(groups) == 0 {
		return nil, nil
	}

	return &User{
		username: username,
		Entry:    entry,
		Groups:   groups,
	}, nil
}
