package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	echo *echo.Echo
	addr string
}

func New(port int, bodyLimit int, corsWhiteList []string) *Server {
	log.Printf("Created rest server with port: %d", port)

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

	return &Server{
		echo: e,
		addr: fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Init() *Server {
	log.Println("Initializing server...")
	s.initApi()
	return s
}

func (s *Server) initApi() {
	log.Println("Initializing api...")

	s.echo.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

}

func (s *Server) Run() error {
	log.Println("Starting server")
	return s.echo.Start(s.addr)
}
