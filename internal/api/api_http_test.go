package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	jwt "github.com/golang-jwt/jwt/v5"
)

func signTestToken(t *testing.T, secret string, userID int) string {
	t.Helper()
	claims := &jwtCustomClaims{
		UserID: userID,
		Role:   "Тест",
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func TestHealth_NoAuth(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestAuthLogin_OK(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{}
	uc.AuthorizeFunc = func(user entities.User) (*bool, error) {
		return boolPtr(true), nil
	}
	uc.SelectUserByEmailFunc = func(email string) (*entities.User, error) {
		return &entities.User{
			ID:        1,
			Email:     email,
			FirstName: "A",
			LastName:  "B",
			Role:      "Админ",
		}, nil
	}

	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	body := bytes.NewBufferString(`{"email":"admin@example.com","password":"Test12345"}`)
	req := httptest.NewRequest(http.MethodPost, "/auth", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d body=%s", rec.Code, rec.Body.String())
	}

	var u entities.User
	if err := json.Unmarshal(rec.Body.Bytes(), &u); err != nil {
		t.Fatalf("unmarshal: %v body=%s", err, rec.Body.String())
	}
	if u.Token == "" {
		t.Fatalf("expected token in response")
	}
}

func TestRoles_RequiresJWT(t *testing.T) {
	t.Parallel()

	uc := &usecaseMock{}
	uc.GetRolesFunc = func() ([]entities.Role, error) {
		return []entities.Role{{ID: 1, Title: "Админ"}}, nil
	}
	s := NewServer("127.0.0.1", 0, uc, "secret", "http://localhost", "./materials")

	// Без токена -> 401 от middleware
	req := httptest.NewRequest(http.MethodGet, "/roles", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}

	// С токеном -> 200
	req2 := httptest.NewRequest(http.MethodGet, "/roles", nil)
	req2.Header.Set("Authorization", "Bearer "+signTestToken(t, "secret", 123))
	rec2 := httptest.NewRecorder()
	s.server.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec2.Code, rec2.Body.String())
	}
}

