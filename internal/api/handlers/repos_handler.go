package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sudzekai-git-api/internal/dal/repositories"
)

type ReposHandler struct {
	repo *repositories.ReposRepo
}

func NewReposHandler(repo *repositories.ReposRepo) *ReposHandler {
	return &ReposHandler{repo: repo}
}

func (h *ReposHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.repo.GetAllRepos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (h *ReposHandler) GetRepoInfo(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/repos/")

	repo, err := h.repo.GetRepoInfo(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)
}

func (h *ReposHandler) GetLastCommits(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/repos/")
	name = strings.TrimSuffix(name, "/commits")

	commits, err := h.repo.GetLastCommits(name, 10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commits)
}

func (h *ReposHandler) CreateRepo(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateRepo(req.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Repository created"})
}

func (h *ReposHandler) DeleteRepo(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/repos/")

	if err := h.repo.DeleteRepo(name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Repository deleted"})
}
