package http

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuth_Success(t *testing.T) {
	jwtSecret := []byte("secret")
	authFunc := NewAuthenticationFunc(jwtSecret)

	// Generate valid token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	// Setup Input
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	input := &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
	}

	// Execute
	err := authFunc(context.Background(), input)

	// Verify
	if err != nil {
		t.Errorf("expected success, got error: %v", err)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	authFunc := NewAuthenticationFunc([]byte("secret"))
	req := httptest.NewRequest("GET", "/", nil)
	input := &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
	}

	err := authFunc(context.Background(), input)
	if err == nil {
		t.Error("expected error for missing header, got nil")
	}
}

func TestAuth_InvalidFormat(t *testing.T) {
	authFunc := NewAuthenticationFunc([]byte("secret"))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidToken")
	input := &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
	}

	err := authFunc(context.Background(), input)
	if err == nil {
		t.Error("expected error for invalid format, got nil")
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	jwtSecret := []byte("secret")
	authFunc := NewAuthenticationFunc(jwtSecret)

	// Generate expired token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	input := &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
	}

	err := authFunc(context.Background(), input)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestAuth_WrongSigningMethod(t *testing.T) {
	jwtSecret := []byte("secret")
	authFunc := NewAuthenticationFunc(jwtSecret)

	// Generate token with different signing method (None)
	token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"sub": "user-123",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	input := &openapi3filter.AuthenticationInput{
		RequestValidationInput: &openapi3filter.RequestValidationInput{
			Request: req,
		},
	}

	err := authFunc(context.Background(), input)
	if err == nil {
		t.Error("expected error for wrong signing method, got nil")
	}
}
