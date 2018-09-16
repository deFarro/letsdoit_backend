package router

import (
	"encoding/json"
	"net/http"

	"github.com/letsdoit_back/app/data"
	"github.com/letsdoit_back/app/database"
)

// HandleTodos return all todos
func HandleTodos(w http.ResponseWriter, r *http.Request) {
	todos := database.FetchTodos()

	payload, err := json.Marshal(todos)
	if err != nil {
		data.SendError("error while marshalling response", w)
	}

	w.Write(payload)
}
