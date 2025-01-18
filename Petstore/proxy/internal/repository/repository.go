package repository

import "context"

type Pet struct {
	Id        int      `json:"id"`
	Category  Category `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Tags      []Tag    `json:"tags"`
	Status    string   `json:"status"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	Id       int    `json:"id"`
	PetId    int    `json:"petId"`
	Quantity int    `json:"quantity"`
	ShipDate string `json:"shipDate"`
	Complete bool   `json:"complete"`
}

var Inventory []Store

var Pets []Pet

type User struct {
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, limit, offset int) ([]User, error)
}
