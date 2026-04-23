package api

import (
	"os"

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

	// Установка уровня логирования на ERROR
	api.server.Logger.SetLevel(log.ERROR)

	// Настройка логирования в файл
	logFile, err := os.OpenFile("logs/api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	api.server.Logger.SetOutput(logFile)

	// Логирование запросов в формате: [время] | статус | метод | IP | URI | время обработки
	api.server.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `[${time_custom}] | ${status} | ${method} | ${remote_ip}${path} |${uri} | ${latency_human}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
		Output:           logFile,
	}))

	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{frontAddress},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	api.address = ip + ":" + string(rune(port))

	return &api
}

func (s *Server) Run() {
	s.server.Logger.Fatal(s.server.Start(s.address))
}
