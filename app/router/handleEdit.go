package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/todo"
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
		SendError("error while reading request payload", w)
		return
	}
	defer r.Body.Close()

	var todo todo.Todo
	err = json.Unmarshal(reqPayload, &todo)
	if err != nil {
		SendError("error while unmarshalling request payload", w)
		return
	}

	switch r.Method {
	case "PUT":
		updatedTodo, err := todo.AddModify(router.Database, sessionID)
		if err != nil {
			SendError(err.Error(), w)
		}

		payload, err := json.Marshal(updatedTodo)
		if err != nil {
			SendError("error while marshalling response", w)
			return		}

		w.Write(payload)

	case "DELETE":
		err = todo.Flush(router.Database, sessionID)
		if err != nil {
			SendError(err.Error(), w)
		}

	default:
		return
	}
}
