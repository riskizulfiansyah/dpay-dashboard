package http

import (
	"context"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
)

// NewAuthenticationFunc returns a middleware function that extracts the JWT token
// from the Authorization header and validates it.
// It matches the signing method (HS256) and secret used in the auth module.
func NewAuthenticationFunc(jwtSecret []byte) func(context.Context, *openapi3filter.AuthenticationInput) error {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		// Extract Bearer token from header
		authHeader := input.RequestValidationInput.Request.Header.Get("Authorization")
		if authHeader == "" {
			return fmt.Errorf("missing authorization header")
		}

		// Check for "Bearer " prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fmt.Errorf("invalid authorization header format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return fmt.Errorf("invalid or expired token")
		}

		// Extract claims (validation only for now)
		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fmt.Errorf("invalid token claims")
		}

		return nil
	}
}
