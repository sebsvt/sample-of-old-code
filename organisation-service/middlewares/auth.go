package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthRequired returns a middleware that validates JWT tokens using the provided secret.
func AuthRequired(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from the request header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		// Extract the token string without the Bearer prefix
		tokenString := authHeader[len(bearerPrefix):]

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(secret), nil
		})

		// Handle errors in token parsing
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Extract claims from the token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, exists := claims["sub"].(string)
			if !exists {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "user ID not found in token",
				})
			}

			// Set the user ID in the context for handlers
			c.Locals("user_id", userID)
			return c.Next()
		}

		// If token is invalid or missing claims
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token claims",
		})
	}
}
