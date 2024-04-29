package handlers

import (
	"api/cmd/internal"
	"net/http"
)

type testPost struct {
	Test    bool   `json:"test"`
	Message string `json:"message,omitempty"`
}

func Ping(w http.ResponseWriter, r *http.Request) {
	internal.Respond(w, http.StatusOK, "PONG!")
}

func Err(w http.ResponseWriter, r *http.Request) {
	internal.RespondErr(w, http.StatusConflict, "resource conflict")
}

func Post(w http.ResponseWriter, r *http.Request) {
	var incoming testPost
	decoded, err := internal.Decode(r.Body, &incoming)
	if err != nil {
		internal.RespondErr(w, http.StatusBadRequest, "bad request. Check request body")
		return
	}
	internal.Respond(w, http.StatusCreated, decoded)
}
