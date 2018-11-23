package session

type Session struct {
	ID string
	UserID string
}

type SessionTransporter interface {
	InsertSession(Session) error
	DropSession(Session) error
}

// DropSession function to clear session
func(s Session) Drop(tr SessionTransporter) error {
	return tr.DropSession(s)
}