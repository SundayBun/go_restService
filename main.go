package main

import (
	"go.uber.org/zap"
	"goWebService/config"
	"goWebService/repository"
	"goWebService/server"
	"log"
	"os"
)

// @title Go REST API
// @version 1.0
// @description Golang REST API

// @contact.name SundayBun
// @contact.url https://github.com/SundayBun
// @contact.email @gmail.com

// @BasePath /api/v1

func main() {
	log.Println("Starting api server")

	cfg := initConfig()

	psqlDB, err := repository.NewPsqlDB(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v", psqlDB.Stats())
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	sugar := logger.Sugar()

	serv := server.NewServer(cfg, psqlDB, sugar)
	if err = serv.Run(); err != nil {
		log.Fatal(err)
	}
	defer psqlDB.Close()
}

func initConfig() *config.Config {

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
