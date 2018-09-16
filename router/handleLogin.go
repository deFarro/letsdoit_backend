package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/letsdoit_back/app/data"
	"github.com/letsdoit_back/app/database"
)

// HandleLogin handles user login request
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	reqPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		data.SendError("error while reading request payload", w)
		return
	}
	r.Body.Close()

	var user data.User
	err = json.Unmarshal(reqPayload, &user)
	if err != nil {
		data.SendError("error while unmarshalling request payload", w)
		return
	}

	dbuser, err := database.FetchUser(user)
	if err != nil {
		data.SendError(err.Error(), w)
		return
	}

	resPayload, err := dbuser.MarshallJSON()
	if err != nil {
		data.SendError("error while marshalling response", w)
		return
	}

	w.Write(resPayload)
}
