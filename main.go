// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"sync/atomic"

// 	"github.com/joho/godotenv"
// 	_ "github.com/lib/pq"
// 	"github.com/yujen77300/Chirpy-Server/internal/database"
// )

// type apiConfig struct {
// 	fileserverHits atomic.Int32
// 	db             *database.Queries
// 	platform       string
// 	jwtSecret      string
// 	polkaKey       string
// }

// func main() {

// 	godotenv.Load()
// 	dbURL := os.Getenv("DB_URL")
// 	if dbURL == "" {
// 		log.Fatal("DB_URL must be set")
// 	}
// 	platform := os.Getenv("PLATFORM")
// 	if platform == "" {
// 		log.Fatal("PLATFORM must be set")
// 	}
// 	jwtSecret := os.Getenv("SECRET")
// 	if jwtSecret == "" {
// 		log.Fatal("JWT_SECRET environment variable is not set")
// 	}
// 	polkaKey := os.Getenv("POLKA_KEY")
// 	if polkaKey == "" {
// 		log.Fatal("POLKA_KEY environment variable is not set")
// 	}

// 	db, err := sql.Open("postgres", dbURL)
// 	if err != nil {
// 		log.Fatalf("Error opening database: %s", err)
// 	}
// 	dbQueries := database.New(db)

// 	cfg := &apiConfig{
// 		db:        dbQueries,
// 		platform:  platform,
// 		jwtSecret: jwtSecret,
// 		polkaKey:  polkaKey,
// 	}
// 	mux := http.NewServeMux()

// 	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

// 	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets")))))

// 	mux.HandleFunc("GET /api/healthz", healthzHandler)
// 	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
// 	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
// 	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
// 	mux.HandleFunc("POST /api/login", cfg.loginHanlder)
// 	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
// 	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
// 	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpByIDHandler)

// 	mux.HandleFunc("POST /api/refresh", cfg.refreshHandler)
// 	mux.HandleFunc("POST /api/revoke", cfg.revokeHanlder)

// 	mux.HandleFunc("PUT /api/users", cfg.updateUserHandler)
// 	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirpByIDHandler)

// 	mux.HandleFunc("POST /api/polka/webhooks", cfg.polkaWebhooksHandler)

// 	server := &http.Server{
// 		Addr:    ":8080",
// 		Handler: mux,
// 	}

// 	fmt.Println("Starting server on :8080")
// 	if err := server.ListenAndServe(); err != nil {
// 		fmt.Println("Server failed:", err)
// 	}
// }

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yujen77300/Chirpy-Server/internal/api"
	"github.com/yujen77300/Chirpy-Server/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: failed to load .env file: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwtSecret := os.Getenv("SECRET")
	if jwtSecret == "" {
		log.Fatal("SECRET environment variable is not set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(db)

	var hits atomic.Int32
	server := api.NewServer(api.ServerConfig{
		DB:             dbQueries,
		Platform:       platform,
		JWTSecret:      jwtSecret,
		PolkaKey:       polkaKey,
		FileserverHits: &hits,
	})

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server.Router()); err != nil {
		fmt.Println("Server failed:", err)
	}
}
