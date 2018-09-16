package router

import (
	"net/http"

	"github.com/letsdoit_back/app/data"
	"github.com/letsdoit_back/app/database"
)

// HandleLogout handles user login request
func HandleLogout(w http.ResponseWriter, r *http.Request) {
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
