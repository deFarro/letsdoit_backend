package data

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Todos is the type for all todos
type Todos []Todo

// Todo is the type for a single todo
type Todo struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	ID          string     `json:"id"`
	Author      PublicUser `json:"author"`
}

// GenerateTodoID generates unique ID for a todo
func (t Todo) GenerateTodoID() string {
	salt := strconv.FormatInt(time.Now().Unix(), 10)

	return fmt.Sprintf("%x", md5.Sum([]byte(t.Title+t.Description+salt)))
}

// User is the type for a user
type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	SessionID    string `json:"sessionID"`
	PasswordHash string `json:"password"`
}

// PublicUser is the type for a user to be passed to client
type PublicUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Public returns public fields for a User
func (u User) Public() PublicUser {
	return PublicUser{
		ID:       u.ID,
		Username: u.Username,
	}
}

// GenerateSessionID generates unique sessionID
func (u User) GenerateSessionID() string {
	salt := strconv.FormatInt(time.Now().Unix(), 10)

	return fmt.Sprintf("%x", md5.Sum([]byte(u.Username+salt)))
}

// MarshallJSON custom json marshaller to hide some fields
func (u *User) MarshallJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		SessionID string `json:"sessionID"`
	}{
		ID:        u.ID,
		Username:  u.Username,
		SessionID: u.SessionID,
	})
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
