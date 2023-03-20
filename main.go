package main

import (
	"github.com/labstack/echo/v4"
	"goWebService/config"
	"goWebService/handler"
	"goWebService/http"
	"goWebService/repository"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting api server")

	cfg := initConfig()

	psqlDB, err := repository.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v", psqlDB.Stats())
	}

	e := echo.New()
	e.Server.ReadTimeout = time.Second * cfg.Server.ReadTimeout
	e.Server.WriteTimeout = time.Second * cfg.Server.WriteTimeout

	v1 := e.Group("/api/v1")
	group := v1.Group("/account")

	http.AccountRoutes(group, handler.NewAccountHandler(cfg, repository.NewPgRepository(psqlDB)))

	go func() {
		log.Printf("Server is listening on PORT: %s", cfg.Server.Port)
		e.Server.ReadTimeout = time.Second * cfg.Server.ReadTimeout
		e.Server.WriteTimeout = time.Second * cfg.Server.WriteTimeout
		//e.Server.MaxHeaderBytes = maxHeaderBytes
		if err := e.Start(cfg.Server.Port); err != nil {
			log.Fatalf("Error starting HTTP Server: ", err)
		}
	}()

	//go func() {
	//	log.Printf("Starting Server on PORT: %s", cfg.Server.PprofPort)
	//	if err := http.ListenAndServe(cfg.Server.PprofPort, http.DefaultServeMux); err != nil {
	//		log.Printf("Error PPROF ListenAndServe: %s", err)
	//	}
	//}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	defer psqlDB.Close()
}

func initConfig() *config.Config {
	os.Setenv("config", "local")

	configPath := GetConfigPath(os.Getenv("config"))

	cfgFile, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}
	return cfg
}

func GetConfigPath(configPath string) string {
	if configPath == "docker" {
		return "./config/config-docker"
	}
	return "./config/config-local"
}
