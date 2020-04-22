package session

//Session holds HTTP session data
type Session interface {
	Username() string
	DisplayName() string
}

//Store is a session storage mechanism
type Store interface {
	//Create creates and returns a session id for the given session
	//or an error if one occurred
	Create(s Session) (id string, err error)
	//Read returns the session for the given id or nil if it doesn't exist
	//or an error if one occurred
	Read(id string) (Session, error)
}
