package api

import (
    "net/http"
    "sync/atomic"

    "github.com/yujen77300/Chirpy-Server/internal/api/handlers"
    "github.com/yujen77300/Chirpy-Server/internal/api/middlewares"
    "github.com/yujen77300/Chirpy-Server/internal/database"
)

type ServerConfig struct {
    DB             *database.Queries
    Platform       string
    JWTSecret      string
    PolkaKey       string
    FileserverHits *atomic.Int32
}

type Server struct {
    config ServerConfig
}

func NewServer(cfg ServerConfig) *Server {
    return &Server{
        config: cfg,
    }
}

// Router sets up the HTTP routes
func (s *Server) Router() http.Handler {
    healthHandler := handlers.NewHealthHandler()
    authHandler := handlers.NewAuthHandler(s.config.DB, s.config.JWTSecret)
    chirpsHandler := handlers.NewChirpsHandler(s.config.DB, s.config.JWTSecret)
    usersHandler := handlers.NewUserHandler(s.config.DB, s.config.JWTSecret)
    adminHandler := handlers.NewAdminHandler(s.config.DB, s.config.Platform, s.config.FileserverHits)
    webhookHandler := handlers.NewWebhookHandler(s.config.DB, s.config.PolkaKey)
    metricsMiddleware := middlewares.NewMetricsMiddleware(s.config.FileserverHits)

    mux := http.NewServeMux()

    mux.Handle("/app/", metricsMiddleware.MetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
    mux.Handle("/app/assets/", metricsMiddleware.MetricsInc(http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets")))))

    mux.HandleFunc("GET /api/healthz", healthHandler.HealthCheck)
    mux.HandleFunc("GET /admin/metrics", adminHandler.GetMetrics)
    mux.HandleFunc("POST /admin/reset", adminHandler.Reset)
    mux.HandleFunc("POST /api/users", usersHandler.Create)
    mux.HandleFunc("PUT /api/users", usersHandler.Update)
    mux.HandleFunc("POST /api/login", authHandler.Login)
    mux.HandleFunc("POST /api/refresh", authHandler.RefreshToken)
    mux.HandleFunc("POST /api/revoke", authHandler.RevokeToken)
    mux.HandleFunc("POST /api/chirps", chirpsHandler.Create)
    mux.HandleFunc("GET /api/chirps", chirpsHandler.GetAll)
    mux.HandleFunc("GET /api/chirps/{chirpID}", chirpsHandler.GetByID)
    mux.HandleFunc("DELETE /api/chirps/{chirpID}", chirpsHandler.Delete)
    mux.HandleFunc("POST /api/polka/webhooks", webhookHandler.HandlePolkaWebhooks)

    return mux
}