package rest

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	v1 "github.com/lapkomo2018/goTwitterAuthService/internal/transport/rest/v1"
	echoSwagger "github.com/swaggo/echo-swagger"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lapkomo2018/goTwitterAuthService/docs"
)

type (
	Config struct {
		Port           int
		BodyLimit      int
		AllowedOrigins []string
	}

	Server struct {
		echo *echo.Echo
		addr string
	}
)

func New(config Config) *Server {
	log.Printf("Creating rest server with port: %d", config.Port)

	e := echo.New()

	e.Use(middleware.BodyLimit(strconv.Itoa(config.BodyLimit)))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} | ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	corsConfig := middleware.CORSConfig{
		AllowOrigins: config.AllowedOrigins,
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return &Server{
		echo: e,
		addr: fmt.Sprintf(":%d", config.Port),
	}
}

func (s *Server) Init(userService v1.UserService, tokenService v1.TokenService, authenticationService v1.AuthenticationService, validator v1.Validator) *Server {
	log.Println("Initializing server...")
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	log.Println("Initializing api...")
	handlerV1 := v1.New(userService, tokenService, authenticationService, validator)
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
