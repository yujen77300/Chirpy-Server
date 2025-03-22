package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/yujen77300/Chirpy-Server/internal/auth"
	"github.com/yujen77300/Chirpy-Server/internal/database"
	"github.com/yujen77300/Chirpy-Server/internal/utils"
)

type WebhookHandler struct {
	db       *database.Queries
	polkaKey string
}

func NewWebhookHandler(db *database.Queries, polkaKey string) *WebhookHandler {
	return &WebhookHandler{
		db:       db,
		polkaKey: polkaKey,
	}
}

func (h *WebhookHandler) HandlePolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type webhook struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Missing or malformed token")
		return
	}

	if apikey != h.polkaKey {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid API key")
		return
	}

	params := webhook{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		utils.RespondWithError(w, http.StatusNoContent, "Invalid event")
		return
	}

	_, err = h.db.UpgradeUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			utils.RespondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Couldn't update user")
		return
	}

	utils.RespondWithJSON(w, http.StatusNoContent, nil)

}
