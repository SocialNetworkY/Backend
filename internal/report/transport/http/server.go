package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SocialNetworkY/Backend/internal/report/transport/http/v1"

	"github.com/SocialNetworkY/Backend/pkg/binder"
	"github.com/SocialNetworkY/Backend/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

func (s *Server) Init(rs v1.ReportService, ag v1.AuthGateway, ug v1.UserGateway) *Server {
	log.Println("Initializing server...")
	log.Println("Initializing api...")
	handlerV1 := v1.New(rs, ag, ug)
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
