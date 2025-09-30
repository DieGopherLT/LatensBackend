package middleware

import (
	"strings"

	"github.com/DieGopherLT/LatensBackend/internal/services/token"
	"github.com/gofiber/fiber/v2"
)

// Boilerplate for custom fiber middleware
func Guard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "No token provided.",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format.",
			})
		}

		tokenString := parts[1]

		// Parse JWT using token service
		payload, err := token.Parse(tokenString)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token.",
			})
		}

		c.Locals("user", *payload)
		return c.Next()
	}
}