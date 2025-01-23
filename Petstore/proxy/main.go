package main

import (
	"Petstore/internal/controller"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

// @title Swagger Petstore
// @version 1.0
// @description Это сваггер магазина питомцев
// @host localhost:8080
// @BasePath
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @RequestAddressSearch представляет запрос для поиска
// @Description Этот эндпоинт позволяет получить адрес по наименованию
// @Param address body ResponseAddress true "Географические координаты"

func main() {
	r := chi.NewRouter()

	// Pet
	r.Post("/pet", controller.AddPetHandler)
	r.Post("/pet/{petId}/uploadImage", controller.UploadImageHandler)
	r.Put("/pet/", controller.UpdatePetHandler)
	r.Get("/pet/findByStatus", controller.FindByStatus)
	r.Get("/pet/{petId}", controller.GetPetByID)
	r.Post("/pet/{petId}", controller.GetPostHandler)
	r.Delete("/pet/{petId}", controller.DeletePetHandler)

	// Store
	r.Post("/store/order", controller.StoreOrder)
	r.Get("/store/inventory", controller.StoreInvent)
	r.Get("/store/order/{orderId}", controller.GetOrderId)
	r.Delete("/store/order/{orderId}", controller.DeleteOrderHandler)

	// User
	r.Post("/user/createWithList", controller.CreateListHandler)
	r.Get("/user/{username}", controller.GettingUsername)
	r.Put("/user/{username}", controller.UpdateUsername)
	r.Delete("/user/{username}", controller.Deleteuser)
	r.Get("/user/login", controller.LogUser)
	r.Get("/user/logout", controller.LogoutUser)
	r.Post("/user/createWithArray", controller.CreateWithArray)
	r.Post("/user", controller.CreateUser)

	fmt.Println("Starting server...")
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
	}
}
