package service

import (
	"Petstore/internal/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func AddPet(pet repository.Pet) (repository.Pet, error) {
	repository.Pets = append(repository.Pets, pet)
	return pet, nil
}

func AddImage(petId string, additionalMetadata string) repository.UploadResponse {

	response := repository.UploadResponse{
		Code:    200,
		Type:    "succes",
		Message: fmt.Sprintf("Image uploaded for pet ID %s with additional metadata: %s", petId, additionalMetadata),
	}
	return response
}

func UpdatePet(w http.ResponseWriter, r http.Request) {
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

func FindStatus(w http.ResponseWriter, r http.Request, status string) {
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

func GetIdPet(w http.ResponseWriter, r http.Request, Id int) {
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
