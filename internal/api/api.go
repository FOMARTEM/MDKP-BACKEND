package api

import (
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
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

type jwtCustomClaims struct {
	UserID int    `json:"id"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
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
		Skipper: func(c echo.Context) bool {
			// Браузерные CORS preflight запросы (OPTIONS) — нормальны и могут засорять лог.
			return c.Request().Method == echo.OPTIONS
		},
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
	normalizedFront := strings.TrimSpace(frontAddress)
	normalizedFront = strings.TrimRight(normalizedFront, "/")
	allowOrigins := []string{"*"}
	if normalizedFront != "" {
		allowOrigins = []string{normalizedFront}
	}
	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  allowOrigins,
		AllowMethods:  []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders: []string{"Content-Disposition"},
	}))

	// JWT middleware для защиты маршрутов, требующих аутентификации
	api.server.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: "header:Authorization:Bearer ",
		ErrorHandler: func(c echo.Context, err error) error {
			return c.JSON(401, map[string]string{"error": "Unauthorized"})
		},
		// Исключаем из JWT middleware маршруты, которые не требуют аутентификации
		Skipper: func(c echo.Context) bool {
			// Разрешаем доступ к маршруту /health без JWT
			if c.Path() == "/health" {
				return true
			}
			// Разрешаем доступ к маршрутам аутентификации без JWT
			if c.Path() == "/auth" || c.Path() == "/auth/refresh" || c.Path() == "/auth/logout" {
				return true
			}

			return false
		},
		// В ContextKey middleware кладёт токен целиком. Дальше в SuccessHandler раскладываем нужные поля в контекст.
		ContextKey: "user",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SuccessHandler: func(c echo.Context) {
			token, ok := c.Get("user").(*jwt.Token)
			if !ok || token == nil {
				return
			}
			claims, ok := token.Claims.(*jwtCustomClaims)
			if !ok || claims == nil {
				return
			}
			// Совместимость: часть хендлеров ожидает ключ "id".
			c.Set("id", claims.UserID)
			c.Set("userID", claims.UserID)
			c.Set("role", claims.Role)
			c.Set("email", claims.Email)
			c.Set("name", claims.Name)
		},
	}))

	api.address = ip + ":" + strconv.Itoa(port)

	api.registerRoutes()

	return &api
}

func (s *Server) Run() {
	s.server.Logger.Fatal(s.server.Start(s.address))
}
