package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/data"
	"github.com/deFarro/letsdoit_backend/app/database"
)

// HandleEdit handles user login request
func HandleEdit(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		DeleteTodo(w, r)

	case "PUT":
		return

	default:
		return
	}
}

// DeleteTodo handles todo removing
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
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

	err = database.FlushTodo(sessionID, todo.ID)
	if err != nil {
		data.SendError(err.Error(), w)
	}
}
