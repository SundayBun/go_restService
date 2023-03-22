package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"goWebService/config"
	_ "goWebService/docs"
	"goWebService/handler"
	"goWebService/http"
	"goWebService/middleware"
	"goWebService/pkg"
	"goWebService/repository"
	"log"
	http3 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg    *config.Config
	echo   *echo.Echo
	db     *sqlx.DB
	logger *zap.SugaredLogger
}

func NewServer(cfg *config.Config, db *sqlx.DB, logger *zap.SugaredLogger) *Server {
	return &Server{cfg: cfg, echo: echo.New(), db: db, logger: logger}
}

func (s Server) Run() error {

	mw := middleware.NewMiddlewareManager(s.cfg, s.logger)

	s.echo.Use(mw.RequestLoggerMiddleware)

	metric := s.createMetrics()
	s.echo.Use(mw.MetricsMiddleware(metric))

	s.echo.GET("/swagger/*", echoSwagger.WrapHandler)

	server := &http3.Server{

		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		log.Printf("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			log.Fatalf("Error starting HTTP Server: ", err)
		}
	}()

	v1 := s.echo.Group("/api/v1")
	group := v1.Group("/account")
	http.AccountRoutes(group, handler.NewAccountHandler(s.cfg, repository.NewPgRepository(s.db)))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	log.Printf("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}

func (s Server) createMetrics() pkg.Metrics {
	metrics, err := pkg.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Info(
		"Metrics available URL: %s, ServiceName: %s",
		s.cfg.Metrics.URL,
		s.cfg.Metrics.ServiceName,
	)
	return metrics
}

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)
