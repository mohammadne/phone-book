package http

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func (s *Server) fetchUserId(c *fiber.Ctx) error {
	headerBytes := c.Request().Header.Peek("Authorization")
	header := strings.TrimPrefix(string(headerBytes), "Bearer ")

	if len(header) == 0 {
		s.logger.Error("Missing authorization header")
		response := "please provide your authentication information"
		return c.Status(http.StatusUnauthorized).SendString(response)
	}

	var id uint64
	if err := s.token.ExtractTokenData(string(header), &id); err != nil || id == 0 {
		s.logger.Error("Invalid token header", zap.Error(err))
		response := "invalid token header, please login again"
		return c.Status(http.StatusBadRequest).SendString(response)
	}

	c.Locals("user-id", id)
	return c.Next()
}
