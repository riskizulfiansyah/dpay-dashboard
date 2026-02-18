package usecase_test

import (
	"testing"
	"time"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/auth/usecase"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Mock Repo
type mockUserRepo struct {
	user *entity.User
	err  error
}

func (m *mockUserRepo) GetUserByEmail(email string) (*entity.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func TestLogin_Success(t *testing.T) {
	// Setup
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		Role:         "admin",
		PasswordHash: string(hashedPassword),
	}

	repo := &mockUserRepo{user: user}
	secret := []byte("secret")
	uc := usecase.NewAuthUsecase(repo, secret, time.Hour)

	// Execute
	tokenString, returnedUser, err := uc.Login(user.Email, password)

	// Verify
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if returnedUser != user {
		t.Errorf("expected user %v, got %v", user, returnedUser)
	}

	// Verify Token
	token, _ := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if !token.Valid {
		t.Error("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("invalid claims")
	}

	if claims["sub"] != user.ID {
		t.Errorf("expected sub %s, got %v", user.ID, claims["sub"])
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{err: entity.ErrorNotFound("user not found")}
	uc := usecase.NewAuthUsecase(repo, []byte("secret"), time.Hour)

	_, _, err := uc.Login("unknown@example.com", "password")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_EmptyUserID(t *testing.T) {
	// Repo returns a user struct but with empty ID (edge case based on implementation)
	repo := &mockUserRepo{user: &entity.User{}}
	uc := usecase.NewAuthUsecase(repo, []byte("secret"), time.Hour)

	_, _, err := uc.Login("test@example.com", "password")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	user := &entity.User{
		ID:           "user-123",
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
	}

	repo := &mockUserRepo{user: user}
	uc := usecase.NewAuthUsecase(repo, []byte("secret"), time.Hour)

	_, _, err := uc.Login("test@example.com", "wrong-password")

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr, ok := err.(*entity.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != entity.ErrorCodeUnauthorized {
		t.Errorf("expected unauthorized code, got %s", appErr.Code)
	}
}
