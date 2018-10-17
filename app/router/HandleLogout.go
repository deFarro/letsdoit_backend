package router

import (
	"net/http"

)

// HandleLogout handles user login request
func (router *Router) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	query := r.URL.Query()
	sessionID := query.Get("sessionID")

	err := router.Database.DropSession(sessionID)
	if err != nil {
		SendError(err.Error(), w)
	}
}
