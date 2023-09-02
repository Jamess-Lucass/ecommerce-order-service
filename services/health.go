package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type HealthService struct {
	db *mongo.Database
}

func NewHealthService(db *mongo.Database) *HealthService {
	return &HealthService{
		db: db,
	}
}

func (s *HealthService) Ping(ctx context.Context) error {
	return s.db.Client().Ping(ctx, nil)
}
