package requests

import "github.com/google/uuid"

type CreateOrderRequest struct {
	Address     string    `json:"address"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	PhoneNumber string    `json:"phoneNumber"`
	UserId      uuid.UUID `json:"userId"`
	Basket      Basket    `json:"basket"`
}

type Basket struct {
	ID    uuid.UUID    `json:"id"`
	Items []BasketItem `json:"items"`
}

type BasketItem struct {
	ID        uuid.UUID `json:"id"`
	CatalogId uuid.UUID `json:"catalogId"`
	Price     float32   `json:"price"`
	Quantity  uint      `json:"quantity"`
}
