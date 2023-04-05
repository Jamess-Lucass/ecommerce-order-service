package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserId      uuid.UUID          `bson:"userId" json:"userId"`
	Address     string             `bson:"address" json:"address"`
	Email       string             `bson:"email" json:"email"`
	Name        string             `bson:"name" json:"name"`
	PhoneNumber string             `bson:"phoneNumber" json:"phoneNumber"`
	Items       []OrderItem        `bson:"items" json:"items"`
}

type OrderItem struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	CatalogId uuid.UUID          `bson:"catalogId" json:"catalogId"`
	Price     float32            `bson:"price" json:"price"`
	Quantity  uint               `bson:"quantity" json:"quantity"`
}
