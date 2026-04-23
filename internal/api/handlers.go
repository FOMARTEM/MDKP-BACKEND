package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthCheckHandler(c echo.Context) error {
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

func (s *Server) authLoginHandler(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) authRefreshHandler(c echo.Context) error      { return s.notImplemented(c) }
func (s *Server) authLogoutHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) accountUpdateHandler(c echo.Context) error    { return s.notImplemented(c) }
func (s *Server) accountGetHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) activityLogListHandler(c echo.Context) error  { return s.notImplemented(c) }
func (s *Server) userCreateHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) userDeleteHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) userRoleUpdateHandler(c echo.Context) error   { return s.notImplemented(c) }
func (s *Server) rolesListHandler(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) statsGetHandler(c echo.Context) error         { return s.notImplemented(c) }
func (s *Server) findUserHandler(c echo.Context) error         { return s.notImplemented(c) }
func (s *Server) taskCreateHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) taskDeleteHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) taskAssignHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) editCreateHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) editStatusUpdateHandler(c echo.Context) error { return s.notImplemented(c) }
func (s *Server) taskGetHandler(c echo.Context) error          { return s.notImplemented(c) }
func (s *Server) tasksListHandler(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) tasksSearchHandler(c echo.Context) error      { return s.notImplemented(c) }
func (s *Server) taskStatusUpdateHandler(c echo.Context) error { return s.notImplemented(c) }
func (s *Server) editsListHandler(c echo.Context) error        { return s.notImplemented(c) }
func (s *Server) materialGetHandler(c echo.Context) error      { return s.notImplemented(c) }
func (s *Server) versionsListHandler(c echo.Context) error     { return s.notImplemented(c) }
func (s *Server) editDeleteHandler(c echo.Context) error       { return s.notImplemented(c) }
func (s *Server) materialUploadHandler(c echo.Context) error   { return s.notImplemented(c) }
func (s *Server) materialDeleteHandler(c echo.Context) error   { return s.notImplemented(c) }
func (s *Server) versionCreateHandler(c echo.Context) error    { return s.notImplemented(c) }
