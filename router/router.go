package router

import (
	"encoding/json"
	"net/http"

	"github.com/letsdoit_back/app/database"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	data := database.FetchTodos()

	payload, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte("error while marshalling response"))
	}

	w.Write(payload)
}