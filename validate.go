package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (cfg *apiConfig) validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Body string `json:"body"`
	}

	var in input
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(in.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := replaceProfaneWords(in.Body)

	response := map[string]string{"cleaned_body": cleanedBody}
	respondWithJSON(w, http.StatusOK, response)
}

func replaceProfaneWords(body string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"

	words := strings.Split(body, " ")
	for i, word := range words {
		if slices.Contains(profaneWords, strings.ToLower(word)) {
			words[i] = replacement
		}
	}

	return strings.Join(words, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
