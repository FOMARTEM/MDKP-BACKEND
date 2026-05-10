package api

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"github.com/FOMARTEM/MDKP-BACKEND/internal/entities"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// Регистрация доступных маршрутов
func (s *Server) registerRoutes() {
	// Здесь регистрируем все маршруты и их обработчики
	s.server.GET("/health", s.healthCheck)

	// Группировка маршрутов для удобства и логической организации
	// Маршруты для аутентификации
	auth := s.server.Group("/auth") // +
	auth.POST("", s.authLogin)      //++
	//auth.POST("/refresh", s.authRefresh)
	//auth.POST("/logout", s.authLogout)

	// Маршруты для управления аккаунтом
	account := s.server.Group("/account")             // ++
	account.PUT("/password", s.accountPasswordUpdate) // ++
	account.GET("/my", s.accountGet)                  // ++

	// Маршрут для получения логов активности
	s.server.GET("/activitylog", s.activityLogList) // ++

	// Маршруты для управления пользователями
	user := s.server.Group("/user")
	user.GET("/list", s.usersList)
	user.POST("", s.userCreate)         // ++
	user.PUT("/active", s.userActive)   // ++
	user.PUT("/role", s.userRoleUpdate) // ++

	// Маршруты для получения ролей, статуса и поиска пользователей
	s.server.GET("/roles", s.rolesList)  // ++
	s.server.GET("/status", s.statusGet) // ++
	//s.server.GET("/finduser", s.findUser) // ++
	s.server.POST("/finduser", s.findUser) // ++

	// Маршруты для управления задачами
	task := s.server.Group("/task")
	task.POST("", s.taskCreate)       // ++
	task.DELETE("/:id", s.taskDelete) // ++
	//task.PUT("/:id/assign", s.taskAssign)       //
	task.GET("/:id", s.taskGet)                 // ++
	task.PUT("/:id/status", s.taskStatusUpdate) // ++

	// Маршруты для получения и поиска задач
	task.GET("/list", s.tasksList) // ++
	//task.GET("/search", s.tasksSearch) //

	// Маршруты для управления правками
	edit := s.server.Group("/edit")
	edit.POST("/:id", s.editCreate)             // ++
	edit.PUT("/:id/status", s.editStatusUpdate) // ++
	//edit.DELETE("/:id", s.editDelete)           //
	edit.GET("/list/:id", s.editsList) // ++
	edit.GET("/:id", s.editGet)        // ++

	// Маршруты для управления материалами
	material := s.server.Group("/material")
	material.GET("/:id", s.materialGet)               // ++
	material.GET("/download/:id", s.materialDownload) // ++
	material.POST("/:id", s.materialUpload)           // ++
	material.GET("/list/:id", s.materialsList)        // ++
	// material.DELETE("/:id", s.materialDelete)      // Удаление файла

	// Маршруты для получения версий и создания новой версии
	version := s.server.Group("/version")
	version.GET("/list/:id", s.versionsList) // ++
	version.POST("/:id", s.versionCreate)    // ++
	version.GET("/:id", s.versionGet)        // ++
	//version.PUT("/:id", s.versionUpdate)   // Обновление версии
}

// Получение limit и offser из query параметра
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

// Получение ID пользователя из токена
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

// Сохранение материала на сервере
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

// Заглушка
func notImplemented(e echo.Context) error {
	return e.JSON(http.StatusNotImplemented, map[string]any{
		"error":  "not implemented",
		"method": e.Request().Method,
		"path":   e.Path(),
	})
}
