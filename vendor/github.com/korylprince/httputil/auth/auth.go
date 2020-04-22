package auth

import "github.com/korylprince/httputil/session"

//Auth represents an authentication mechanism
type Auth interface {
	//Authenticate authenticates the given credentials and returns a Session initialized with the user's data
	//if the credentials are valid, or nil if not. If an error occurs it is returned.
	Authenticate(username, password string) (session.Session, error)
}
