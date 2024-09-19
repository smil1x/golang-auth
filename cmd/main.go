package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"golang-auth/internal/app/apiserver"
	"golang-auth/internal/app/handler"
	"golang-auth/internal/app/repository"
	"golang-auth/internal/app/service"
	"log"
)

func main() {
	config, err := apiserver.LoadConfig()
	if err != nil {
		log.Fatalf("invalid server configuration: %s", err.Error())
	}

	db, err := repository.ConnectPostgresDB(config.DB)
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	rep := repository.NewRepository(db)
	services := service.NewService(rep, &config.JWT)
	handlers := handler.NewHandlers(services)

	fmt.Println("server started")
	var srv = apiserver.New(config)
	err = srv.Run(handlers.InitHandlers())
	if err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}
