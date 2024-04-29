package internal

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

type errorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type successResponse struct {
	Success bool       `json:"success"`
	Data    json.Token `json:"data,omitempty"`
}

func send(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("COntent-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}

func RespondErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Sending 5xx ERROR response -", msg)
	}

	res, err := json.Marshal(errorResponse{Success: false, Error: msg})
	if err != nil {
		log.Println("Failed to marshal ERROR response -", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	send(w, code, res)
}

func Respond(w http.ResponseWriter, code int, data json.Token) {
	res, err := json.Marshal(successResponse{Success: true, Data: data})
	if err != nil {
		log.Println("Failed to marshal SUCCESS response -", err)
		RespondErr(w, http.StatusInternalServerError, "unexpected error. Try again later")
		return
	}
	send(w, code, res)
}

func Decode[T any](reader io.Reader, dest *T) (*T, error) {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&dest); err != nil {
		log.Println("Failed to decoded INCOMING request body -", err)
		return nil, err
	}

	// ensure only single item in req body
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		log.Println("Multiple items in INCOMING request body -", err)
		return nil, err
	}

	return dest, nil
}
