package http

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/http/v1"
	"github.com/lapkomo2018/goTwitterServices/pkg/binder"
	"github.com/lapkomo2018/goTwitterServices/pkg/validator"
	"log"
	"net/http"
	"time"
)

type Server struct {
	echo *echo.Echo
	addr string
}

func New(bodyLimit string, allowedOrigins []string, port int) *Server {
	log.Printf("Creating rest server with port: %d", port)

	e := echo.New()
	e.Binder = binder.NewEchoCustomBinder()
	e.Validator = validator.NewEchoCustomValidator()
	e.Use(middleware.Recover())

	e.Use(middleware.BodyLimit(bodyLimit))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} | ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	corsConfig := middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.Pre(middleware.RemoveTrailingSlash())

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return &Server{
		echo: e,
		addr: fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Init(userService v1.UserService, tokenService v1.TokenService, authenticationService v1.AuthenticationService, refreshTokenDuration time.Duration) *Server {
	log.Println("Initializing server...")
	log.Println("Initializing api...")
	handlerV1 := v1.New(userService, tokenService, authenticationService, refreshTokenDuration)
	api := s.echo.Group("/api")
	{
		handlerV1.Init(api)
	}

	return s
}

func (s *Server) Run() error {
	log.Println("Starting server")
	return s.echo.Start(s.addr)
}
