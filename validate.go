package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Body string `json:"body"`
	}

	var in input
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(`{"error": "Something went wrong"}`))
		return
	}

	if len(in.Body) > 140 {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte(`{"error": "Chirp is too long"}`))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(`{"valid": true}`))
}
