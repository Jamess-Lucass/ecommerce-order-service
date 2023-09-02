package handlers

import "github.com/gofiber/fiber/v2"

func (s *Server) Healthz(c *fiber.Ctx) error {

	if err := s.healthService.Ping(c.Context()); err != nil {
		s.logger.Sugar().Errorf("error occured while pinging redis: %v", err)
		return c.Status(fiber.StatusServiceUnavailable).Send([]byte("Unhealthy"))
	}

	return c.Status(fiber.StatusOK).Send([]byte("Healthy"))
}
