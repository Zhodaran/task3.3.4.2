package control

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
)

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

type CreateResponse struct {
	Message string `json:"message"`
}

type rErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

// @Summary Create SQL user
// @Description This description created new SQL user
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.User true "User login details"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/users [post]
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uc.userRepo.Create(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Create successful"})
	w.WriteHeader(http.StatusCreated)
}

// @Summary Get SQL user
// @Description This description created new SQL user
// @Tags Controller
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} CreateResponse "Greate successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/users/{id} [get]
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user, err := uc.userRepo.GetByID(context.Background(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Great successful"})
	json.NewEncoder(w).Encode(user)
}

// @Summary Update SQL user
// @Description This description created new SQL user
// @Tags Controller
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body repository.User true "User detail to update"
// @Success 200 {object} CreateResponse "Update successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/users{id} [put]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := uc.userRepo.Update(context.Background(), user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Update successful"})
	w.WriteHeader(http.StatusNoContent)
}

// @Summary Delete SQL user
// @Description This description created new SQL user
// @Tags Controller
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} CreateResponse "Delete successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/users/{id} [delete]
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := uc.userRepo.Delete(context.Background(), id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "Delete successful"})
	w.WriteHeader(http.StatusNoContent)
}

// @Summary List SQL user
// @Description This description created new SQL user
// @Tags Controller
// @Accept json
// @Produce json
// @Success 200 {object} CreateResponse "List successful"
// @Failure 400 {object} rErrorResponse "Invalid request"
// @Failure 401 {object} rErrorResponse "Invalid credentials"
// @Failure 500 {object} rErrorResponse "Internal server error"
// @Router /api/users [get]
func (uc *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit := 10 // Установите значение по умолчанию
	offset := 0 // Установите значение по умолчанию
	users, err := uc.userRepo.List(context.Background(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(CreateResponse{Message: "List successful"})
	json.NewEncoder(w).Encode(users)
}
