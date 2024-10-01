package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lapkomo2018/goTwitterServices/internal/auth/transport/rest/v1"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	Config struct {
		BodyLimit      int
		AllowedOrigins []string
	}

	Server struct {
		echo *echo.Echo
		addr string
	}
)

func New(cfg Config, port int) *Server {
	log.Printf("Creating rest server with port: %d", port)

	e := echo.New()

	e.Use(middleware.BodyLimit(strconv.Itoa(cfg.BodyLimit)))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} | ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	corsConfig := middleware.CORSConfig{
		AllowOrigins: cfg.AllowedOrigins,
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

func (s *Server) Init(userService v1.UserService, tokenService v1.TokenService, authenticationService v1.AuthenticationService, validator v1.Validator, refreshTokenDuration time.Duration) *Server {
	log.Println("Initializing server...")
	log.Println("Initializing api...")
	handlerV1 := v1.New(userService, tokenService, authenticationService, validator, refreshTokenDuration)
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
