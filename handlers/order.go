package handlers

import (
	"github.com/Jamess-Lucass/ecommerce-order-service/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET: /api/v1/orders
func (s *Server) GetOrders(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claim)

	orders, err := s.orderService.List(c.Context(), claims)
	if err != nil {
		s.logger.Sugar().Infof("error getting orders: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: "Could not load orders."})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

// GET: /api/v1/orders/1
func (s *Server) GetOrder(c *fiber.Ctx) error {
	id, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Code: fiber.StatusBadRequest, Message: err.Error()})
	}
	claims := c.Locals("claims").(*middleware.Claim)

	order, err := s.orderService.Get(c.Context(), claims, id)
	if err != nil {
		s.logger.Sugar().Infof("error getting order: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{Code: fiber.StatusNotFound, Message: "Could not find order."})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}
