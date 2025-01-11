package repository

import "context"

// User представляет собой модель пользователя
type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	DeletedAt *string `json:"deleted_at"` // Для логического удаления
}

// UserRepository определяет методы для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]User, error)
}
