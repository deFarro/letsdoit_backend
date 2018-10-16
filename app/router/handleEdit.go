package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/data"
	"github.com/deFarro/letsdoit_backend/app/database"
)

// HandleEdit handles user login request
func (router *Router) HandleEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	query := r.URL.Query()
	sessionID := query.Get("sessionID")

	reqPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data.SendError("error while reading request payload", w)
		return
	}
	defer r.Body.Close()

	var todo data.Todo
	err = json.Unmarshal(reqPayload, &todo)
	if err != nil {
		data.SendError("error while unmarshalling request payload", w)
		return
	}

	switch r.Method {
	case "PUT":
		updatedTodo, err := database.UpdateTodo(sessionID, todo)
		if err != nil {
			data.SendError(err.Error(), w)
		}

		payload, err := json.Marshal(updatedTodo)
		if err != nil {
			data.SendError("error while marshalling response", w)
			return
		}

		w.Write(payload)

	case "DELETE":
		err = database.FlushTodo(sessionID, todo.ID)
		if err != nil {
			data.SendError(err.Error(), w)
		}

	default:
		return
	}
}
