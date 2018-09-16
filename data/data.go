package data

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Todos is the type for all todos
type Todos []Todo

// Todo is the type for a single todo
type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          string `json:"id"`
	Author      User   `json:"author"`
}

// User is the type for a user
type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
	Password  string `json:"password"`
}

// Error type for errors
type Error struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

// NewError generates new error
func NewError(m string) Error {
	return Error{
		Error:   true,
		Message: m,
	}
}

// SendError sends error to client
func SendError(m string, w http.ResponseWriter) {
	err := NewError(m)

	payload, _ := json.Marshal(err)

	w.Write(payload)
}

// GenerateSessionID generates unique sessionID
func (u User) GenerateSessionID() string {
	salt := string(time.Now().Unix())

	return fmt.Sprintf("%x", md5.Sum([]byte(u.Username+salt)))
}

// MarshallJSON custom json marshaller to hide some fields
func (u *User) MarshallJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		SessionID string `json:"sessionId"`
	}{
		ID:        u.ID,
		Username:  u.Username,
		SessionID: u.SessionID,
	})
}
