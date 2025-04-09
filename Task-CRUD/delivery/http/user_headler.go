package http

import (
	"Task-CRUD/internal/entity"
	"Task-CRUD/internal/usecase"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUC usecase.UserUseCaseInterface
}

func NewUserHandler(userUC usecase.UserUseCaseInterface) *UserHandler {
	return &UserHandler{userUC: userUC}
}

// --- Helper ---

func writeUserError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// --- Handlers ---

// GET /users
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	users, err := h.userUC.GetUsers(ctx)
	if err != nil {
		log.Printf("ERROR | GetUsers: %v", err)
		writeUserError(w, http.StatusInternalServerError, "Gagal mengambil data user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GET /users/{id}
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		writeUserError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	ctx := r.Context()
	user, err := h.userUC.GetUserByID(ctx, id)
	if err != nil {
		log.Printf("ERROR | GetUserByID: %v", err)
		writeUserError(w, http.StatusNotFound, "User tidak ditemukan")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// POST /users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeUserError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if user.Name == "" || user.Email == "" {
		writeUserError(w, http.StatusBadRequest, "Nama dan Email harus diisi")
		return
	}

	log.Printf("DEBUG | CreateUser payload: %+v", user)

	ctx := r.Context()
	if err := h.userUC.CreateUser(ctx, &user); err != nil {
		log.Printf("ERROR | CreateUser: %v", err)
		writeUserError(w, http.StatusInternalServerError, "Gagal membuat user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil dibuat"})
}

// PUT /users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		writeUserError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var user entity.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeUserError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if user.Name == "" {
		writeUserError(w, http.StatusBadRequest, "Nama tidak boleh kosong")
		return
	}

	ctx := r.Context()
	if err := h.userUC.UpdateUser(ctx, id, &user); err != nil {
		log.Printf("ERROR | UpdateUser: %v", err)
		writeUserError(w, http.StatusInternalServerError, "Gagal memperbarui user")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User berhasil diperbarui"})
}

// DELETE /users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromVars(r)
	if err != nil {
		writeUserError(w, http.StatusBadRequest, "ID tidak valid")
		return
	}

	ctx := r.Context()
	if err := h.userUC.DeleteUser(ctx, id); err != nil {
		log.Printf("ERROR | DeleteUser: %v", err)
		writeUserError(w, http.StatusInternalServerError, "Gagal menghapus user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Fprint(w, `{"message": "User berhasil dihapus"}`)
}

// --- Utility ---

func parseIDFromVars(r *http.Request) (uint, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID")
	}
	return uint(id), nil
}
