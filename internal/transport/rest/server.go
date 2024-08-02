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

type Server struct {
	echo *echo.Echo
}

func New(bodyLimit int, corsWhiteList []string) *Server {
	log.Printf("Creating rest server")

	e := echo.New()

	e.Use(middleware.BodyLimit(strconv.Itoa(bodyLimit)))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${time_custom} | ${status} | ${latency_human} | ${remote_ip} | ${method} | ${uri} | ${error}\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	corsConfig := middleware.CORSConfig{
		AllowOrigins: corsWhiteList,
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return &Server{
		echo: e,
	}
}

func (s *Server) Init(userService v1.UserService, tokenService v1.TokenService) *Server {
	log.Println("Initializing server...")
	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)
	s.echo.GET("/swagger", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	log.Println("Initializing api...")
	handlerV1 := v1.New(userService, tokenService)
	api := s.echo.Group("/api")
	{
		handlerV1.Init(api)
	}

	return s
}

func (s *Server) Run(port int) error {
	log.Println("Starting server")
	return s.echo.Start(fmt.Sprintf(":%d", port))
}
