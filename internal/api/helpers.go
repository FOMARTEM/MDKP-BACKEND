package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strconv"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (s *Server) getLimitOffset(e echo.Context) (int, int) {
	limit, _ := strconv.Atoi(e.QueryParam("limit"))
	offset, _ := strconv.Atoi(e.QueryParam("offset"))

	if limit <= 0 {
		limit = 32
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func (s *Server) userIDFromToken(e echo.Context) (int, error) {
	if v := e.Get("id"); v != nil {
		if id, ok := v.(int); ok {
			return id, nil
		}
		return 0, echo.NewHTTPError(500, "invalid id type in context")
	}

	if v := e.Get("userID"); v != nil {
		if id, ok := v.(int); ok {
			return id, nil
		}
		return 0, echo.NewHTTPError(500, "invalid userID type in context")
	}

	// Fallback: если SuccessHandler не сработал, пробуем вытащить из самого токена.
	token, ok := e.Get("user").(*jwt.Token)
	if !ok || token == nil {
		return 0, echo.NewHTTPError(401, "missing jwt token")
	}

	switch claims := token.Claims.(type) {
	case *jwtCustomClaims:
		if claims == nil || claims.UserID == 0 {
			return 0, echo.NewHTTPError(401, "invalid jwt claims")
		}
		return claims.UserID, nil
	case jwt.MapClaims:
		raw, ok := claims["id"]
		if !ok {
			return 0, echo.NewHTTPError(401, "missing id in jwt claims")
		}
		switch n := raw.(type) {
		case float64:
			return int(n), nil
		case int:
			return n, nil
		default:
			return 0, echo.NewHTTPError(401, "invalid id type in jwt claims")
		}
	default:
		return 0, echo.NewHTTPError(401, "unsupported jwt claims type")
	}
}

func (s *Server) saveMaterialFile(file *multipart.FileHeader, material entities.Material) error {
	dst := fmt.Sprintf("./materials/%s-%d.%s", material.Title, material.ID, material.Extension)

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	//на случай если нету папки
	os.MkdirAll("./materials", 0755)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
