package session

type Session struct {
	ID string
	UserID string
}

// Resource interface to operate over sessions
type Resource interface {
	InsertSession(Session) error
	DropSession(Session) error
}

// Insert function to add new session
func(s Session) Insert(r Resource) error {
	return r.InsertSession(s)
}

// DropSession function to clear session
func(s Session) Drop(r Resource) error {
	return r.DropSession(s)
}