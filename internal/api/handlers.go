package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{
		"status": "ok",
	})
}

func (s *Server) notImplemented(c echo.Context) error {
	return c.JSON(http.StatusNotImplemented, map[string]any{
		"error":  "not implemented",
		"method": c.Request().Method,
		"path":   c.Path(),
	})
}

func (s *Server) authLogin(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) authRefresh(c echo.Context) error      { return s.notImplemented(c) }
func (s *Server) authLogout(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) accountUpdate(c echo.Context) error    { return s.notImplemented(c) }
func (s *Server) accountGet(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) activityLogList(c echo.Context) error  { return s.notImplemented(c) }
func (s *Server) userCreate(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) userDelete(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) userRoleUpdate(c echo.Context) error   { return s.notImplemented(c) }
func (s *Server) rolesList(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) statsGet(c echo.Context) error         { return s.notImplemented(c) }
func (s *Server) findUser(c echo.Context) error         { return s.notImplemented(c) }
func (s *Server) taskCreate(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) taskDelete(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) taskAssign(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) editCreate(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) editStatusUpdate(c echo.Context) error { return s.notImplemented(c) }
func (s *Server) taskGet(c echo.Context) error          { return s.notImplemented(c) }
func (s *Server) tasksList(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) tasksSearch(c echo.Context) error      { return s.notImplemented(c) }
func (s *Server) taskStatusUpdate(c echo.Context) error { return s.notImplemented(c) }
func (s *Server) editsList(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) materialGet(c echo.Context) error      { return s.notImplemented(c) }
func (s *Server) versionsList(c echo.Context) error     { return s.notImplemented(c) }
func (s *Server) editDelete(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) materialUpload(c echo.Context) error   { return s.notImplemented(c) }
func (s *Server) materialDelete(c echo.Context) error   { return s.notImplemented(c) }
func (s *Server) versionCreate(c echo.Context) error    { return s.notImplemented(c) }
