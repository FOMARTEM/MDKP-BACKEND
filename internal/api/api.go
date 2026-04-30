package api

import (
	"os"
	"strconv"
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
	s.server.GET("/roles", s.rolesList)   // ++
	s.server.GET("/status", s.statusGet)  // ++
	s.server.GET("/finduser", s.findUser) // ++

	// Маршруты для управления задачами
	task := s.server.Group("/task")
	task.POST("", s.taskCreate)                 // ++
	task.DELETE("/:id", s.taskDelete)           // ++
	task.PUT("/:id/assign", s.taskAssign)       //
	task.GET("/:id", s.taskGet)                 // ++
	task.PUT("/:id/status", s.taskStatusUpdate) // ++

	// Маршруты для получения и поиска задач
	task.GET("/list", s.tasksList)
	task.GET("/search", s.tasksSearch)

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

func getLimitOffset(e echo.Context) (int, int) {
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

func userIDFromToken(e echo.Context) (int, error) {
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
