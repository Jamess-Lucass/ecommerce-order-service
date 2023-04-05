package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Claim struct {
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

func parseToken(tokenString string) (*Claim, error) {
	secret, exists := os.LookupEnv("JWT_SECRET")
	claims := &Claim{}

	if !exists {
		panic("JWT Secret not set")
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claim)

	if !ok {
		return nil, fmt.Errorf("ID token is invalid")
	}

	return claims, nil
}

func getToken(c *fiber.Ctx) string {
	// Auth bearer token
	authHeader := c.Get("Authorization")

	if authHeader != "" {
		tokenHeader := strings.Split(authHeader, "Bearer ")

		if len(tokenHeader) == 2 {
			return tokenHeader[1]
		}
	}

	// Cookies
	cookie := c.Cookies("x-access-token")

	return cookie
}

func JWT() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := getToken(c)

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": fiber.StatusUnauthorized, "message": "Unauthorized"})
		}

		claims, err := parseToken(token)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"code": fiber.StatusUnauthorized, "message": "Unauthorized"})
		}

		c.Locals("claims", claims)

		return c.Next()
	}
}
