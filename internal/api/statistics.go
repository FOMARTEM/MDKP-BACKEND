package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Проверка состояния backenda
func (s *Server) healthCheck(e echo.Context) error {
	return e.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}
