package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"

	"github.com/antgngo/go-http-server/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "Forbidden", nil)
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(r.Context())

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	fmt.Println(dbURL)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("error")
	}

	cfg := &apiConfig{
		db:       database.New(db),
		platform: "dev",
	}

	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	handler := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(handler)))
	serveMux.Handle("/app/assets", http.StripPrefix("/app", handler))
	serveMux.HandleFunc("POST /admin/reset", cfg.handlerReset)

	serveMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	serveMux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		numHits := cfg.fileserverHits.Load()
		w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", numHits)))
	})

	serveMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	serveMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, req *http.Request) {
		type email struct {
			Email string `json:"email"`
		}
		decoder := json.NewDecoder(req.Body)
		reqParams := email{}
		err := decoder.Decode(&reqParams)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error", err)
		}
		user, err := cfg.db.CreateUser(req.Context(), reqParams.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
			return
		}

		newUser := User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		}
		respondWithJSON(w, 201, newUser)
	})

	server.ListenAndServe()
}
