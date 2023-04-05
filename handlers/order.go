package handlers

import (
	"github.com/Jamess-Lucass/ecommerce-order-service/middleware"
	"github.com/gofiber/fiber/v2"
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
