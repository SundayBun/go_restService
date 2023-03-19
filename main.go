package main

import (
	"github.com/gin-gonic/gin"
	"goWebService/config"
	"goWebService/handler"
	"goWebService/repository"
	"log"
	"os"
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

	inMemoryStorage := repository.NewInMemoryStorage()
	hnd := handler.NewHandler(inMemoryStorage)

	router := gin.Default()
	router.GET("/account", hnd.GetAccount)
	router.POST("/account/new", hnd.CreateAccount)
	router.POST("/account/update", hnd.UpdateAccount)
	router.DELETE("/account/delete", hnd.DeleteAccount)

	router.Run()

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
