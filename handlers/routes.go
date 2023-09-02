package handlers

import (
	"os"
	"strings"

	"github.com/Jamess-Lucass/ecommerce-order-service/middleware"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ErrorResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

func (s *Server) Start() error {
	f := fiber.New()
	f.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
		AllowOriginsFunc: func(origin string) bool {
			return strings.EqualFold(os.Getenv("ENVIRONMENT"), "development")
		},
		AllowCredentials: true,
		MaxAge:           0,
	}))

	f.Use(fiberzap.New(fiberzap.Config{
		Logger: s.logger,
	}))

	f.Get("/api/healthz", s.Healthz)

	f.Get("/api/v1/orders", middleware.JWT(), s.GetOrders)
	f.Get("/api/v1/orders/:id", middleware.JWT(), s.GetOrder)

	f.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"code": fiber.StatusNotFound, "message": "No resource found"})
	})

	return f.Listen(":8080")
}
