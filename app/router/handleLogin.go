package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/deFarro/letsdoit_backend/app/data"
	"github.com/deFarro/letsdoit_backend/app/database"
)

// HandleLogin handles user login request
func (router *Router) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}

	reqPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		SendError("error while reading request payload", w)
		return
	}
	defer r.Body.Close()

	var user data.User
	err = json.Unmarshal(reqPayload, &user)
	if err != nil {
		SendError("error while unmarshalling request payload", w)
		return
	}

	dbuser, err := database.FetchUser(router.Database, user)
	if err != nil {
		SendError(err.Error(), w)
		return
	}

	resPayload, err := dbuser.MarshallJSON()
	if err != nil {
		SendError("error while marshalling response", w)
		return
	}

	w.Write(resPayload)
}
