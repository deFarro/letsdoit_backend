package router

import (
	"encoding/json"
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/database"
)

// HandleTodos return all todos
func (router *Router) HandleTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	todos, err := database.FetchTodos(router.Database)
	if err != nil {
		SendError(err.Error(), w)
		return
	}

	payload, err := json.Marshal(todos)
	if err != nil {
		SendError("error while marshalling response", w)
		return
	}

	w.Write(payload)
}
