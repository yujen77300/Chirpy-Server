package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/yujen77300/Chirpy-Server/internal/database"

	"github.com/google/uuid"
)

type response struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
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

	userID, err := uuid.Parse(in.UserID)
	if err != nil {
		log.Printf("Invalid user ID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var chirpParam database.CreateChirpParams
	chirpParam.Body = cleanedBody
	chirpParam.UserID = userID

	chirp, err := cfg.db.CreateChirp(r.Context(), chirpParam)
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	var chirpResponses []response
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, response{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, http.StatusOK, chirpResponses)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpId)
	if err != nil {
		log.Printf("Invalid chirp ID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, http.StatusNotFound, "Something went wrong")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
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

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
