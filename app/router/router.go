package router

import (
	"github.com/deFarro/letsdoit_backend/app/config"
	"github.com/deFarro/letsdoit_backend/app/database"
	"net/http"
	"encoding/json"
	"github.com/deFarro/letsdoit_backend/app/user"
	"github.com/deFarro/letsdoit_backend/app/todo"
	"github.com/deFarro/letsdoit_backend/app/session"
)

type Router struct {
	Settings config.Config
	Database database.Database
	UserTransporter user.UserTransporter
	TodoTransporter todo.TodoTransporter
	SessionTransporter session.SessionTransporter
}

// NewRouter creates new router instance
func NewRouter(settings config.Config) (Router, error) {
	db, err := database.NewDatabase(settings)
	if err != nil {
		return Router{}, err
	}

	return Router{
		Settings: settings,
		Database: db,
	}, nil
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
