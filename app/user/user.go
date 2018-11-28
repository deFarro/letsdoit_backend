package user

import (
	"fmt"
	"crypto/md5"
	"errors"
	"strconv"
	"time"
	"encoding/json"
	"github.com/deFarro/letsdoit_backend/app/session"
)

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

// Resource interface to operate over users
type Resource interface {
	GetUserByID(string) (User, error)
	GetUserBySessionID(string) (User, error)
	InsertSession(session.Session) error
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

// FetchUser fetches user from db, checks password, create user session and returns user
func (u User) Fetch(r Resource) (User, error) {
	userID := fmt.Sprintf("%x", md5.Sum([]byte(u.Username)))

	currentUser, err := r.GetUserByID(userID)
	if err != nil {
		return User{}, err
	}

	if currentUser.PasswordHash != u.PasswordHash {
		return User{}, errors.New("wrong password")
	}

	currentUser.SessionID = currentUser.GenerateSessionID()

	return currentUser, nil
}
