package api

import jwt "github.com/golang-jwt/jwt/v5"

// jwtCustomClaims описывает полезную нагрузку JWT, которую мы ждём от фронтенда.
// Используется только для парсинга/валидации входящих токенов в echo-jwt middleware.
type jwtCustomClaims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}
