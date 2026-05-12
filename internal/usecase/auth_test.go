package usecase

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_AndCompare(t *testing.T) {
	t.Parallel()

	hash, err := hashPassword("Test12345")
	if err != nil {
		t.Fatalf("hashPassword: %v", err)
	}
	if hash == "" {
		t.Fatalf("expected non-empty hash")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("Test12345")); err != nil {
		t.Fatalf("compare: %v", err)
	}
}

func TestAuthorize_Success(t *testing.T) {
	t.Parallel()

	hash, err := hashPassword("pw")
	if err != nil {
		t.Fatalf("hashPassword: %v", err)
	}

	var logs []string
	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return &entities.User{ID: 10, Email: email}, nil
		},
		GetUserPasswordHashByEmailFun: func(email string) (string, error) {
			return hash, nil
		},
		CreateLogFunc: func(userID int, action string) error {
			logs = append(logs, action)
			return nil
		},
	}
	u := NewUsecase(p)

	ok, err := u.Authorize(entities.User{Email: "a@b.c", Password: "pw"})
	if err != nil {
		t.Fatalf("Authorize: %v", err)
	}
	if ok == nil || !*ok {
		t.Fatalf("expected ok=true, got %#v", ok)
	}
	if len(logs) != 1 || logs[0] != "Successful login" {
		t.Fatalf("unexpected logs: %#v", logs)
	}
}

func TestAuthorize_InvalidPassword(t *testing.T) {
	t.Parallel()

	hash, err := hashPassword("correct")
	if err != nil {
		t.Fatalf("hashPassword: %v", err)
	}

	var logs []string
	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return &entities.User{ID: 10, Email: email}, nil
		},
		GetUserPasswordHashByEmailFun: func(email string) (string, error) {
			return hash, nil
		},
		CreateLogFunc: func(userID int, action string) error {
			logs = append(logs, action)
			return nil
		},
	}
	u := NewUsecase(p)

	ok, err := u.Authorize(entities.User{Email: "a@b.c", Password: "wrong"})
	if err == nil || !errors.Is(err, entities.ErrInvalidPassword) {
		t.Fatalf("expected ErrInvalidPassword, got %v", err)
	}
	if ok == nil || *ok {
		t.Fatalf("expected ok=false, got %#v", ok)
	}
	if len(logs) != 1 || logs[0] != "Failed login attempt" {
		t.Fatalf("unexpected logs: %#v", logs)
	}
}

func TestAuthorize_UserNotFound_FromGetUserByEmail(t *testing.T) {
	t.Parallel()

	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return nil, errors.New("db down")
		},
	}
	u := NewUsecase(p)

	_, err := u.Authorize(entities.User{Email: "a@b.c", Password: "pw"})
	if err == nil || !errors.Is(err, entities.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestAuthorize_UserNotFound_FromNoRowsHash(t *testing.T) {
	t.Parallel()

	var logs []string
	p := &providerMock{
		GetUserByEmailFunc: func(email string) (*entities.User, error) {
			return &entities.User{ID: 10, Email: email}, nil
		},
		GetUserPasswordHashByEmailFun: func(email string) (string, error) {
			return "", sql.ErrNoRows
		},
		CreateLogFunc: func(userID int, action string) error {
			logs = append(logs, action)
			return nil
		},
	}
	u := NewUsecase(p)

	ok, err := u.Authorize(entities.User{Email: "a@b.c", Password: "pw"})
	if err == nil || !errors.Is(err, entities.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
	if ok == nil || *ok {
		t.Fatalf("expected ok=false, got %#v", ok)
	}
	if len(logs) != 1 || logs[0] != "Failed login attempt - user not found" {
		t.Fatalf("unexpected logs: %#v", logs)
	}
}

