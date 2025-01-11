package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	_ "studentgit.kata.academy/Zhodaran/go-kata/docs"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"studentgit.kata.academy/Zhodaran/go-kata/controller"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/auth"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/control"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/repository"
	"studentgit.kata.academy/Zhodaran/go-kata/internal/service"
)

// @title Address API
// @version 1.0
// @description API для поиска
// @host localhost:8080
// @BasePath
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @RequestAddressSearch представляет запрос для поиска
// @Description Этот эндпоинт позволяет получить адрес по наименованию
// @Param address body ResponseAddress true "Географические координаты"

type GeocodeRequest struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

// TokenResponse представляет ответ с токеном

// LoginResponse представляет ответ при успешном входе

type Server struct {
	http.Server
}

func (s *Server) Serve() {
	log.Println("Starting server...")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: &v", err)
	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))
	fmt.Println("DB_PASSWORD:", os.Getenv("DB_PASSWORD"))
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))
	fmt.Println("DB_PORT:", os.Getenv("DB_PORT"))
	fmt.Println("Запуск задержки")
	time.Sleep(10 * time.Second)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	runMigrations(db)

	userRepo := repository.NewPostgresUserRepository(db)
	userController := control.NewUserController(userRepo)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	geoService := service.NewGeoService("d9e0649452a137b73d941aa4fb4fcac859372c8c", "ec99b849ebf21277ec821c63e1a2bc8221900b1d") // Создаем новый экземпляр GeoService
	resp := controller.NewResponder(logger)
	r := router(userController, resp, geoService)

	srv := &Server{
		Server: http.Server{
			Addr:         ":8080",
			Handler:      r,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}

	go srv.Serve()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при завершении работы: %v\n", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}

func runMigrations(db *sql.DB) {
	migrationSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		email VARCHAR(255) NOT NULL,
		deleted_at TIMESTAMP NULL
	);`

	_, err := db.Exec(migrationSQL)
	if err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}
}

func proxyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			next.ServeHTTP(w, r)
			return
		}
		proxyURL, _ := url.Parse("http://hugo:1313")
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)
		proxy.ServeHTTP(w, r)
	})
}

func TokenAuthMiddleware(resp controller.Responder) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				resp.ErrorUnauthorized(w, errors.New("missing authorization token"))
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			_, err := auth.TokenAuth.Decode(token)
			if err != nil {
				resp.ErrorUnauthorized(w, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func geocodeHandler(resp controller.Responder, geoService service.GeoProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req GeocodeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.ErrorBadRequest(w, err)
			return
		}

		geo, err := geoService.GetGeoCoordinatesGeocode(req.Lat, req.Lng)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		resp.OutputJSON(w, geo)
	}
}

func searchHandler(resp controller.Responder, geoService service.GeoProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddressSearch
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.ErrorBadRequest(w, err)
			return
		}

		geo, err := geoService.GetGeoCoordinatesAddress(req.Query)
		if err != nil {
			resp.ErrorInternal(w, err)
			return
		}

		resp.OutputJSON(w, geo)
	}
}

func router(userController *control.UserController, resp controller.Responder, geoService service.GeoProvider) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(proxyMiddleware)
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Post("/api/register", auth.Register)
	r.Post("/api/login", auth.Login)
	r.Post("/api/users", userController.CreateUser)        // Создание пользователя
	r.Get("/api/users/{id}", userController.GetUser)       // Получение пользователя по ID
	r.Put("/api/users/{id}", userController.UpdateUser)    // Обновление пользователя
	r.Delete("/api/users/{id}", userController.DeleteUser) // Удаление пользователя
	r.Get("/api/users", userController.ListUsers)

	// Используем обработчики с middleware
	r.With(TokenAuthMiddleware(resp)).Post("/api/address/geocode", geocodeHandler(resp, geoService))
	r.With(TokenAuthMiddleware(resp)).Post("/api/address/search", searchHandler(resp, geoService))

	return r
}
