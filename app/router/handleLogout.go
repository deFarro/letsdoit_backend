package router

import (
	"net/http"
	"github.com/deFarro/letsdoit_backend/app/session"
)

// HandleLogout handles user login request
func (router *Router) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	query := r.URL.Query()
	session := session.Session{
		ID: query.Get("sessionID"),
	}

	err := session.Drop(router.Database)
	if err != nil {
		SendError(err.Error(), w)
	}
}
