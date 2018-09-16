package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/letsdoit_back/app/data"
	"github.com/letsdoit_back/app/database"
)

// HandleUser handles user request
func HandleUser(w http.ResponseWriter, r *http.Request) {
	reqPayload, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.Write([]byte("error while reading request payload"))
		return
	}

	var user data.User
	err = json.Unmarshal(reqPayload, &user)
	if err != nil {
		w.Write([]byte("error while unmarshalling request payload"))
		return
	}

	data, err := database.FetchUser(user)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	resPayload, err := data.MarshallJSON()
	if err != nil {
		w.Write([]byte("error while marshalling response"))
		return
	}

	w.Write(resPayload)
}
