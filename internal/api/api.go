package api

import (
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type Server struct {
	server  *echo.Echo
	address string

	materialStaticPath string

	secretKey string

	uc Usecase
}

func NewServer(ip string, port int, uc Usecase, secretKey string, frontAddress string, materialPath string) *Server {
	api := Server{
		uc:                 uc,
		secretKey:          secretKey,
		materialStaticPath: materialPath,
	}

	// Инициализация Echo сервера
	api.server = echo.New()

	// Установка уровня логирования на INFO (нужно для логирования запросов)
	api.server.Logger.SetLevel(log.INFO)

	// Настройка логирования в файл
	logFile, err := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	api.server.Logger.SetOutput(logFile)

	// Логирование запросов в формате: [время] | статус | метод | IP | URI | время обработки
	api.server.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogLatency:  true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			timestamp := v.StartTime.Format(time.RFC3339Nano)
			if v.Error != nil {
				api.server.Logger.Errorf("%s | %d | %s | %s | %s | %v | err=%v", timestamp, v.Status, v.Method, v.RemoteIP, v.URI, v.Latency, v.Error)
				return nil
			}
			api.server.Logger.Infof("%s | %d | %s | %s | %s | %v", timestamp, v.Status, v.Method, v.RemoteIP, v.URI, v.Latency)
			return nil
		},
	}))

	// CORS middleware для разрешения запросов с фронтенда
	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontAddress},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	api.address = ip + ":" + strconv.Itoa(port)

	api.registerRoutes()

	return &api
}

func (s *Server) Run() {
	s.server.Logger.Fatal(s.server.Start(s.address))
}

func (s *Server) registerRoutes() {
	// Здесь регистрируем все маршруты и их обработчики
	s.server.GET("/health", s.healthCheck)

	// Группировка маршрутов для удобства и логической организации
	// Маршруты для аутентификации
	auth := s.server.Group("/auth")
	auth.POST("", s.authLogin)
	auth.POST("/refresh", s.authRefresh)
	auth.POST("/logout", s.authLogout)

	// Маршруты для управления аккаунтом
	account := s.server.Group("/account")
	account.PUT("", s.accountUpdate)
	account.GET("", s.accountGet)

	// Маршрут для получения логов активности
	s.server.GET("/activitylog", s.activityLogList)

	// Маршруты для управления пользователями
	s.server.POST("/user", s.userCreate)
	s.server.DELETE("/user/:id", s.userDelete)
	s.server.PUT("/user/role", s.userRoleUpdate)

	// Маршруты для получения ролей, статистики и поиска пользователей
	s.server.GET("/roles", s.rolesList)
	s.server.GET("/stats", s.statsGet)
	s.server.GET("/finduser", s.findUser)

	// Маршруты для управления задачами
	task := s.server.Group("/task")
	task.POST("", s.taskCreate)
	task.DELETE("/:id", s.taskDelete)
	task.PUT("/:id/assign", s.taskAssign)
	task.GET("/:id", s.taskGet)
	task.PUT("/:id/status", s.taskStatusUpdate)

	// Маршруты для получения и поиска задач
	s.server.GET("/tasks", s.tasksList)
	s.server.GET("/tasks/search", s.tasksSearch)

	// Маршруты для управления правками
	edit := s.server.Group("/edit")
	edit.POST("/:id", s.editCreate)
	edit.PUT("/:id/status", s.editStatusUpdate)
	edit.DELETE("/:id", s.editDelete)

	// Маршрут для получения правок по ID
	s.server.GET("/edits/:id", s.editsList)

	// Маршруты для управления материалами
	material := s.server.Group("/material")
	material.GET("/:id", s.materialGet)
	material.POST("", s.materialUpload)
	material.DELETE("/:id", s.materialDelete)

	// Маршруты для получения версий и создания новой версии
	s.server.GET("/versions/:id", s.versionsList)
	s.server.POST("/version", s.versionCreate)
}
