package main

import (
	"goWebService/config"
	"goWebService/repository"
	"goWebService/server"
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

	serv := server.NewServer(cfg, psqlDB)
	if err = serv.Run(); err != nil {
		log.Fatal(err)
	}
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
