package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func TestGetLimitOffset_Defaults(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s := &Server{}
	limit, offset := s.getLimitOffset(c)
	if limit != 32 || offset != 0 {
		t.Fatalf("unexpected limit/offset: %d/%d", limit, offset)
	}
}

func TestGetLimitOffset_Parse(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/x?limit=10&offset=5", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	s := &Server{}
	limit, offset := s.getLimitOffset(c)
	if limit != 10 || offset != 5 {
		t.Fatalf("unexpected limit/offset: %d/%d", limit, offset)
	}
}

func TestUserIDFromToken_FromContextID(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("id", 7)

	s := &Server{}
	id, err := s.userIDFromToken(c)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 7 {
		t.Fatalf("unexpected id: %d", id)
	}
}

func TestUserIDFromToken_FallbackJWTMapClaims(t *testing.T) {
	t.Parallel()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"id": float64(9)})
	c.Set("user", token)

	s := &Server{}
	id, err := s.userIDFromToken(c)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if id != 9 {
		t.Fatalf("unexpected id: %d", id)
	}
}

