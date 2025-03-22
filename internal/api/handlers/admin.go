package handlers

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/yujen77300/Chirpy-Server/internal/database"
)

type AdminHandler struct {
	db             *database.Queries
	platform       string
	fileserverHits *atomic.Int32
}

func NewAdminHandler(db *database.Queries, platform string, fileserverHits *atomic.Int32) *AdminHandler {
	return &AdminHandler{
		db:             db,
		platform:       platform,
		fileserverHits: fileserverHits,
	}
}



func (h *AdminHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`,
		h.fileserverHits.Load())))
}

func (h *AdminHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if h.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Forbidden: this endpoint is only available in development mode"))
		return
	}

	// Delete all users from the database
	err := h.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("Error deleting users: %s", err)))
		return
	}

	h.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte("Hits reset to 0"))
}
