package repository_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/durianpay/fullstack-boilerplate/internal/entity"
	"github.com/durianpay/fullstack-boilerplate/internal/module/auth/repository"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		t.Fatalf("failed to create users table: %v", err)
	}

	return db
}

func seedUser(t *testing.T, db *sql.DB, user entity.User) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO users (id, email, password_hash, role) VALUES (?, ?, ?, ?)`,
		user.ID, user.Email, user.PasswordHash, user.Role,
	)
	if err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	testUser := entity.User{
		ID:           "test-user-id",
		Email:        "test@example.com",
		PasswordHash: "hashed-secret",
		Role:         "admin",
	}
	seedUser(t, db, testUser)

	repo := repository.NewUserRepo(db)

	foundUser, err := repo.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if foundUser.ID != testUser.ID {
		t.Errorf("expected ID %s, got %s", testUser.ID, foundUser.ID)
	}
	if foundUser.Email != testUser.Email {
		t.Errorf("expected Email %s, got %s", testUser.Email, foundUser.Email)
	}
	if foundUser.Role != testUser.Role {
		t.Errorf("expected Role %s, got %s", testUser.Role, foundUser.Role)
	}
	if foundUser.PasswordHash != testUser.PasswordHash {
		t.Errorf("expected Hash %s, got %s", testUser.PasswordHash, foundUser.PasswordHash)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := repository.NewUserRepo(db)

	_, err := repo.GetUserByEmail("nonexistent@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	appErr, ok := err.(*entity.AppError)
	if !ok {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != entity.ErrorCodeNotFound {
		t.Errorf("expected not_found code, got %s", appErr.Code)
	}
}
