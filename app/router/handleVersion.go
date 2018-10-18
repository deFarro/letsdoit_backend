package router

import (
	"net/http"
	"encoding/json"
)

type VersionResponse struct {
	Version string `json:"version"`
}

// HandleVersion handles user login request
func (router *Router) HandleVersion(w http.ResponseWriter, r *http.Request) {
	resPayload, err := json.Marshal(VersionResponse{
		Version: router.Settings.Version,
	})
	if err != nil {
		SendError("error while marshalling response", w)
		return
	}

	w.Write(resPayload)
}
