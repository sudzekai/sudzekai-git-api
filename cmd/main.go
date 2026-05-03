package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"sudzekai-git-api/internal/api/handlers"
	"sudzekai-git-api/internal/dal/repositories"
)

// CORS middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	var (
		reposPath = flag.String("repos-path", "repos", "Path to git repositories folder")
		apiPort   = flag.String("port", "8081", "API server port")
	)

	flag.Parse()

	// Создаём папку для репозиториев, если её нет
	if err := os.MkdirAll(*reposPath, 0755); err != nil {
		log.Fatalf("Failed to create repos directory: %v", err)
	}

	log.Printf("Git repositories path: %s", *reposPath)
	log.Printf("API server starting on port %s", *apiPort)

	mux := http.NewServeMux()
	RegisterReposEndpoints(mux, *reposPath)

	handler := enableCORS(mux)

	serverAddr := ":" + *apiPort
	log.Printf("Server starting on %s", serverAddr)

	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// RegisterReposEndpoints регистрирует все endpoints
func RegisterReposEndpoints(mux *http.ServeMux, reposPath string) {
	repo := repositories.NewReposRepo(reposPath)
	handler := handlers.NewReposHandler(repo)

	// GET endpoints
	mux.HandleFunc("GET /api/repos", handler.GetAllRepos)
	mux.HandleFunc("GET /api/repos/{name}", handler.GetRepoInfo)
	mux.HandleFunc("GET /api/repos/{name}/commits", handler.GetLastCommits)

	// POST endpoints
	mux.HandleFunc("POST /api/repos", handler.CreateRepo)

	// DELETE endpoints
	mux.HandleFunc("DELETE /api/repos/{name}", handler.DeleteRepo)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
}
