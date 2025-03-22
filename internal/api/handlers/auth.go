package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yujen77300/Chirpy-Server/internal/auth"
	"github.com/yujen77300/Chirpy-Server/internal/database"
	"github.com/yujen77300/Chirpy-Server/internal/models"
	"github.com/yujen77300/Chirpy-Server/internal/utils"
)

type AuthHandler struct {
	db        *database.Queries
	jwtSecret string
}

func NewAuthHandler(db *database.Queries, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		models.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// if params.ExpiresIn == 0 || params.ExpiresIn > 3600 {
	// 	params.ExpiresIn = 3600
	// }

	user, err := h.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, h.jwtSecret, time.Hour)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token")
		return
	}

	_, err = h.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  refreshToken,
		// Refresh tokens should expire after 60 days.
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		User: models.User{
			ID:          user.ID,
			Email:       user.Email,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	user, err := h.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token")
		return
	}

	accessToken, err := auth.MakeJWT(
		user.ID,
		h.jwtSecret,
		time.Hour,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Couldn't validate token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

func (h *AuthHandler) RevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	_, err = h.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
