package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "text/json; charset=utf-8")
	w.WriteHeader(status)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("unable to json encode response:\n%v\nwith error: %v\n", payload, err)
		respondWithError(w, http.StatusInternalServerError, "internal error")
	}

	_, err = w.Write(body)
	if err != nil {
		log.Printf("unable to write response: %v\n", err)
		respondWithError(w, http.StatusInternalServerError, "internal error")
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{
		Error: msg,
	})
}
