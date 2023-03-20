package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"goWebService/config"
	"goWebService/handler"
	"goWebService/http"
	"goWebService/repository"
	"log"
	http3 "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg  *config.Config
	echo *echo.Echo
	db   *sqlx.DB
}

func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	return &Server{cfg: cfg, echo: echo.New(), db: db}
}

func (s Server) Run() error {

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

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)
