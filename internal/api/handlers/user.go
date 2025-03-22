package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/yujen77300/Chirpy-Server/internal/auth"
	"github.com/yujen77300/Chirpy-Server/internal/database"
	"github.com/yujen77300/Chirpy-Server/internal/utils"
		"github.com/yujen77300/Chirpy-Server/internal/models"
)

type UserHandler struct {
	db        *database.Queries
	jwtSecret string
}

func NewUserHandler(db *database.Queries, jwtSecret string) *UserHandler {
	return &UserHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}


func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		models.User
	}

	in := input{}
	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(in.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := h.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          in.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, response{
		User: models.User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})

}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		models.User
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

	in := input{}
	err = json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		log.Printf("Error decoding JSON: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(in.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	user, err := h.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          in.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: models.User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
	})

}
