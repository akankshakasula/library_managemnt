package middleware

import (
	"errors" 
	"fmt"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	"library-management/internal/db"
	"library-management/internal/handlers" 
)

func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header format must be Bearer <token>"})
		}

		tokenString := parts[1]

		claims := &handlers.Claims{} 
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(db.JWTSecret), nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token signature"})
			}
			if errors.Is(err, jwt.ErrTokenExpired) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token expired"})
			}
			if errors.Is(err, jwt.ErrTokenNotValidYet) || errors.Is(err, jwt.ErrTokenMalformed) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format or not valid yet"})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token: " + err.Error()})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		c.Locals("userID", claims.UserID)
		c.Locals("userEmail", claims.Email)
		c.Locals("userRole", claims.Role)
		return c.Next()
	}
}

func Authorize(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User role not found in context"})
		}

		for _, role := range requiredRoles {
			if userRole == role {
				return c.Next() 
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Insufficient permissions."})
	}
}
