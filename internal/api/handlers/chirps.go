package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yujen77300/Chirpy-Server/internal/auth"
	"github.com/yujen77300/Chirpy-Server/internal/database"
	"github.com/yujen77300/Chirpy-Server/internal/utils"
)

type ChirpsHandler struct {
	db        *database.Queries
	jwtSecret string
}

func NewChirpsHandler(db *database.Queries, jwtSecret string) *ChirpsHandler {
	return &ChirpsHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

type chirpResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

// Create handles the creation of new chirps
func (h *ChirpsHandler) Create(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing or malformed token")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, h.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	var params struct {
		Body string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) > 140 {
		utils.RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleanedBody := replaceProfaneWords(params.Body)

	chirp, err := h.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})

	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

// GetAll returns all chirps, with optional filtering
func (h *ChirpsHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	authorIDStr := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	var chirps []database.Chirp
	var err error

	if authorIDStr != "" {
		authorID, parseErr := uuid.Parse(authorIDStr)
		if parseErr != nil {
			log.Printf("Invalid author ID: %s", err)
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}

		chirps, err = h.db.GetChirpsByAuthorID(r.Context(), authorID)
	} else {
		chirps, err = h.db.GetChirps(r.Context())
	}

	if err != nil {
		log.Printf("Error getting chirps: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}


	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	} else {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	}

	var chirpResponses []chirpResponse
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, chirpResponse{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, chirpResponses)
}

// GetByID returns a single chirp by ID
func (h *ChirpsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpIDStr)
	if err != nil {
		log.Printf("Invalid chirp ID: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	chirp, err := h.db.GetChirp(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, chirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

// Delete removes a chirp if the user is the author
func (h *ChirpsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")

	id, err := uuid.Parse(chirpIDStr)
	if err != nil {
		log.Printf("Invalid chirp ID: %s", err)
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing or malformed token")
		return
	}

	userID, err := auth.ValidateJWT(tokenString, h.jwtSecret)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	chirp, err := h.db.GetChirp(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}

	if chirp.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "You cannot delete another user's chirp")
		return
	}

	err = h.db.DeleteChirp(r.Context(), chirp.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)
}

// Helper function for profanity filtering
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
