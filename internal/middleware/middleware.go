package middleware

import (
	"strings"

	"github.com/DieGopherLT/LatensBackend/internal/services/token"
	"github.com/gofiber/fiber/v2"
)

// Boilerplate for custom fiber middleware
func Guard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authToken := c.Cookies("auth_token")

		parts := strings.Split(authToken, " ")
		tokenType, tokenString := parts[0], parts[1]
		if len(parts) != 2 || tokenType != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token format.",
			})
		}

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