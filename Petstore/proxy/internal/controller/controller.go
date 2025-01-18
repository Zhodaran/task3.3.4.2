package controller

import (
	"Petstore/internal/repository"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
)

type UploadResponse struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type CreateResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	BadRequest      string `json:"400"`
	DadataBad       string `json:"500"`
	SuccefulRequest string `json:"200"`
}

// @Summary Add pet handler
// @Description This description addadder new pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet [post]
func AddPetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var pet repository.Pet

	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		return
	}
	repository.Pets = append(repository.Pets, pet)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pet)
}

// @Summary Download image pet
// @Description This description upload image pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet/{petId}/uploadImage [post]
func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	petId := chi.URLParam(r, "petId")
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}

	additionalMetadata := r.FormValue("additionalMetadata")

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Unable to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	response := UploadResponse{
		Code:    200,
		Type:    "succes",
		Message: fmt.Sprintf("Image uploaded for pet ID %s with additional metadata: %s", petId, additionalMetadata),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// @Summary Update pet
// @Description This description update pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet/ [put]
func UpdatePetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var pet repository.Pet

	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pet); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary Find pet
// @Description This description update pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet/findByStatus [get]
func FindByStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "Status paramete if required", http.StatusBadRequest)
		return
	}

	var filteredPets []repository.Pet
	for _, pet := range repository.Pets {
		if strings.EqualFold(pet.Status, status) {
			filteredPets = append(filteredPets, pet)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(filteredPets); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary Update pet
// @Description This description update pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet/{petId} [get]
func GetPetByID(w http.ResponseWriter, r *http.Request) {
	// Получаем petId из URL
	id := chi.URLParam(r, "petId")
	Id, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Error reverse")
		return
	}
	var foundPet *repository.Pet
	for _, pet := range repository.Pets {
		if pet.Id == Id {
			foundPet = &pet
			break
		}
	}

	if foundPet == nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем питомца в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(foundPet); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary Update pet
// @Description This description update pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet/{petId} [post]
func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем petId из URL
	id := chi.URLParam(r, "petId")

	// Преобразуем строку в целое число
	petId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid pet ID", http.StatusBadRequest)
		return
	}

	// Получаем данные из формы
	name := r.FormValue("name")
	status := r.FormValue("status")

	// Ищем питомца в "базе данных"
	var foundPet *repository.Pet
	for i, pet := range repository.Pets {
		if pet.Id == petId {
			// Обновляем данные питомца
			repository.Pets[i].Name = name
			repository.Pets[i].Status = status
			foundPet = &repository.Pets[i]
			break
		}
	}

	if foundPet == nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем обновленного питомца в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(foundPet); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary Update pet
// @Description This description update pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param user body repository.Pet true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /pet/{petId} [delete]
func DeletePetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем petId из URL
	id := chi.URLParam(r, "petId")

	// Преобразуем строку в целое число
	petId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID supplied", http.StatusBadRequest)
		return
	}

	// Ищем питомца в "базе данных"
	var foundIndex int
	var found bool
	for i, pet := range repository.Pets {
		if pet.Id == petId {
			foundIndex = i
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}

	// Удаляем питомца из "базы данных"
	repository.Pets = append(repository.Pets[:foundIndex], repository.Pets[foundIndex+1:]...)

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем успешный ответ
	response := map[string]string{"message": "Pet deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func StoreInvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("incorect method")
		return
	}
	var inventory *[]repository.Store

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем питомца в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(inventory); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
