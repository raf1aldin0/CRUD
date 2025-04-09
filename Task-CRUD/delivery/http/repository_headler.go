package http

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/usecase"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type RepoHandler struct {
	repoUC usecase.RepoUseCaseInterface
}

func NewRepoHandler(repoUC usecase.RepoUseCaseInterface) *RepoHandler {
	return &RepoHandler{repoUC: repoUC}
}

// --- Helper ---

func writeRepoError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func parseRepoID(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil || id == 0 {
		return 0, err
	}
	return uint(id), nil
}

// --- Handlers ---

// GET /repositories
func (h *RepoHandler) GetAllRepos(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	repos, err := h.repoUC.GetAllRepos(ctx)
	if err != nil {
		log.Printf("ERROR | GetAllRepos: %v", err)
		writeRepoError(w, http.StatusInternalServerError, "Gagal mengambil daftar repository")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

// GET /repositories/{id}
func (h *RepoHandler) GetRepositoryByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseRepoID(r)
	if err != nil {
		writeRepoError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	ctx := r.Context()
	repo, err := h.repoUC.GetRepositoryByID(ctx, id)
	if err != nil {
		log.Printf("ERROR | GetRepositoryByID: %v", err)
		writeRepoError(w, http.StatusNotFound, "Repository tidak ditemukan")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repo)
}

// POST /repositories
func (h *RepoHandler) CreateRepo(w http.ResponseWriter, r *http.Request) {
	var repo entity.Repository
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&repo); err != nil {
		writeRepoError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if repo.Name == "" || repo.UserID == 0 {
		writeRepoError(w, http.StatusBadRequest, "Nama repository dan user_id tidak boleh kosong")
		return
	}

	ctx := r.Context()
	if err := h.repoUC.CreateRepo(ctx, &repo); err != nil {
		log.Printf("ERROR | CreateRepo: %v", err)
		writeRepoError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(repo)
}

// PUT /repositories/{id}
func (h *RepoHandler) UpdateRepo(w http.ResponseWriter, r *http.Request) {
	id, err := parseRepoID(r)
	if err != nil {
		writeRepoError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var updatedRepo entity.Repository
	if err := json.NewDecoder(r.Body).Decode(&updatedRepo); err != nil {
		writeRepoError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	ctx := r.Context()
	if err := h.repoUC.UpdateRepo(ctx, id, &updatedRepo); err != nil {
		log.Printf("ERROR | UpdateRepo: %v", err)
		writeRepoError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedRepo)
}

// DELETE /repositories/{id}
func (h *RepoHandler) DeleteRepo(w http.ResponseWriter, r *http.Request) {
	id, err := parseRepoID(r)
	if err != nil {
		writeRepoError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	ctx := r.Context()
	if err := h.repoUC.DeleteRepo(ctx, id); err != nil {
		log.Printf("ERROR | DeleteRepo: %v", err)
		writeRepoError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
