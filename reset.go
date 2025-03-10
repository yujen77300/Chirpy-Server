package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Forbidden: this endpoint is only available in development mode"))
		return
	}

	// Delete all users from the database
	err := cfg.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Error deleting users: %s", err)))
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Hits reset to 0"))
}
