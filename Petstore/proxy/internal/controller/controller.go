package controller

import (
	"Petstore/internal/repository"
	"Petstore/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

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
		http.Error(w, "FAil", http.StatusBadRequest)
	}

	createdPet, err := service.AddPet(pet)
	if err != nil {
		http.Error(w, "Fail", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdPet)
}

// @Summary Download image pet
// @Description This description upload image pet
// @Tags Controller
// @Accept multipart/from-data
// @Produce json
// @Param petId path int64 true "file to addadder"
// @Param additionalMetadata formData string false "Additional data to pass tp server"
// @Param file formData file true "File to upload"
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
	response := service.AddImage(petId, additionalMetadata)

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
	service.UpdatePet(w, *r)

}

// @Summary Find pet
// @Description This description update pet
// @Tags Controller
// @Accept json
// @Produce json
// @Param status query string true "Status values that need to be considered for filter" Ennums(availabe, pending, sold)
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
	service.FindStatus(w, *r, status)

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
	service.GetIdPet(w, *r, Id)
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

// @Summary Id order
// @Description This description update pet
// @Tags Store
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /store/order [get]
func StoreInvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Print("incorect method")
		return
	}
	statusCount := make(map[string]int)
	for _, order := range repository.Orders {
		statusCount[order.Status]++
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем питомца в JSON и отправляем ответ
	if err := json.NewEncoder(w).Encode(statusCount); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary Id order
// @Description This description update pet
// @Tags Store
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /store/order [post]
func StoreOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println("Error Method Post")
		return
	}
	var order repository.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	order.ShipDate = time.Now().UTC().Format(time.RFC3339Nano)

	repository.Orders[order.Id] = order

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary Id order
// @Description This description update pet
// @Tags Store
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /store/order/{orderId} [get]
func GetOrderId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Println("Error Method Post")
		return
	}

	id := chi.URLParam(r, "orderId")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID supplied", http.StatusBadRequest)
		return
	}
	var foundOrder *repository.Order
	for _, order := range repository.Orders {
		if order.Id == orderId {
			foundOrder = &order
			break
		}
	}

	if foundOrder == nil {
		http.Error(w, "Pet not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(foundOrder); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary Delete order
// @Description This description update pet
// @Tags Store
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /store/order/{orderId} [delete]
func DeleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем petId из URL
	id := chi.URLParam(r, "orderId")

	// Преобразуем строку в целое число
	orderId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID supplied", http.StatusBadRequest)
		return
	}

	// Ищем питомца в "базе данных"
	var foundIndex int
	var found bool
	for i, order := range repository.Orders {
		if order.Id == orderId {
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
	delete(repository.Orders, foundIndex)
	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем успешный ответ
	response := map[string]string{"message": "Pet deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/createWithList [post]
func CreateListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method now allowed", http.StatusBadRequest)
		return
	}
	var users []repository.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/{username} [get]
func GettingUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Lol this is not this method", http.StatusBadRequest)
		return
	}
	username := chi.URLParam(r, "username")
	var foundUser repository.User
	for _, user := range repository.Users {
		if username == user.Username {
			foundUser = user
			return
		}
	}
	if err := json.NewEncoder(w).Encode(foundUser); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/{username} [put]
func UpdateUsername(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var user repository.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/{username} [delete]
func Deleteuser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Not pu method", http.StatusBadRequest)
		return
	}
	username := chi.URLParam(r, "username")
	var found bool = false
	var foundIndex int
	for i, user := range repository.Users {
		if username == user.Username {
			found = true
			foundIndex = i
			return
		}
	}
	if !found {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	delete(repository.Users, foundIndex)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Отправляем успешный ответ
	response := map[string]string{"message": "User deleted successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/login [get]
func LogUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username == "" || password == "" {
		http.Error(w, "Invalid username/password supplied", http.StatusBadRequest)
		return
	}
	var foundId int
	var found bool = false
	for i, user := range repository.Users {
		if username == user.Username {
			if password == user.Password {
				foundId = i
				found = true
				return
			}
		}
	}
	if !found {
		http.Error(w, "Error authorizaztion", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message": "User logged in successfully",
		"userId":  foundId,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/logout [get]
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Http error", http.StatusRequestEntityTooLarge)
		return
	}

	response := repository.UploadResponse{
		Code:    200,
		Type:    "unknown",
		Message: "ok",
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем ответ в JSON и отправляем его клиенту
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/createWithArray [post]
func CreateWithArray(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем входящие данные
	var users []repository.User
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	// Здесь можно добавить логику для сохранения пользователей в базу данных

	// Формируем ответ
	response := repository.UploadResponse{
		Code:    200,
		Type:    "unknown",
		Message: "ok",
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем ответ в JSON и отправляем его клиенту
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}

// @Summary List user order
// @Description This description update pet
// @Tags user
// @Accept json
// @Produce json
// @Param user body repository.Order true "Pet addadder"
// @Success 200 {object} CreateResponse "Create successful"
// @Failure 400 {object} ErrorResponse "Invalid request"
// @Failure 401 {object} ErrorResponse "Invalid credentials"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /user/user [post]
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Декодируем входящие данные
	var user repository.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	// Здесь можно добавить логику для сохранения пользователя в базу данных

	// Формируем ответ
	response := repository.UploadResponse{
		Code:    200,
		Type:    "unknown",
		Message: "User created successfully",
	}

	// Устанавливаем заголовок ответа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Кодируем ответ в JSON и отправляем его клиенту
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка при формировании ответа", http.StatusInternalServerError)
	}
}
