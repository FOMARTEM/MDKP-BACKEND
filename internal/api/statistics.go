package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthCheck(e echo.Context) error {
	return e.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func notImplemented(e echo.Context) error {
	return e.JSON(http.StatusNotImplemented, map[string]any{
		"error":  "not implemented",
		"method": e.Request().Method,
		"path":   e.Path(),
	})
}
