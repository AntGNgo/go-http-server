package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func main() {
	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	cfg := &apiConfig{}

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

	server.ListenAndServe()
}
