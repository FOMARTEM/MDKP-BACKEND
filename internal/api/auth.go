package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Аутентификация и получение токена JWT ++
func (s *Server) authLogin(e echo.Context) error {
	var user entities.User

	err := e.Bind(&user)
	if err != nil {
		return e.JSON(http.StatusBadRequest, err.Error())
	}

	login, err := s.uc.Authorize(user)
	if err != nil {
		return e.JSON(http.StatusUnauthorized, err.Error())
	}

	if !*login {
		return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	u, err := s.uc.SelectUserByEmail(user.Email)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	claims := jwtCustomClaims{
		UserID: u.ID,
		Role:   u.Role,
		Email:  u.Email,
		Name:   fmt.Sprintf("%s %s", u.LastName, u.FirstName),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 8000)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	u.Token, err = token.SignedString([]byte(s.secretKey))

	if err != nil {
		return e.JSON(http.StatusInternalServerError, err.Error())
	}

	return e.JSON(http.StatusOK, u)
}

//func (s *Server) authRefresh(e echo.Context) error { return notImplemented(e) }
//func (s *Server) authLogout(e echo.Context) error  { return notImplemented(e) }
