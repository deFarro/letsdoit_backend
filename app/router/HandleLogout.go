package router

import (
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/data"
	"github.com/deFarro/letsdoit_backend/app/database"
)

// HandleLogout handles user login request
func (router *Router) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	query := r.URL.Query()
	sessionID := query.Get("sessionID")

	err := database.DropSession(sessionID)
	if err != nil {
		data.SendError(err.Error(), w)
	}
}
