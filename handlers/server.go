package handlers

import (
	"github.com/Jamess-Lucass/ecommerce-order-service/services"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Server struct {
	validator    *validator.Validate
	logger       *zap.Logger
	orderService *services.OrderService
}

func NewServer(logger *zap.Logger, orderService *services.OrderService) *Server {
	return &Server{
		validator:    validator.New(),
		logger:       logger,
		orderService: orderService,
	}
}
