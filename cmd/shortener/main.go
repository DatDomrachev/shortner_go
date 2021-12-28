package main

import (
	"context"
	"github.com/DatDomrachev/shortner_go/internal/app/config"
	"github.com/DatDomrachev/shortner_go/internal/app/repository"
	"github.com/DatDomrachev/shortner_go/internal/app/database"
	"github.com/DatDomrachev/shortner_go/internal/app/server"
	"log"
	"os"
	"os/signal"
)

func main() {

	config, err := config.New()

	if err != nil {
		log.Fatalf("failed to configurate:+%v", err)
	}

	config.InitFlags()
	db, err := database.New(config.DBURL) 
	repo := repository.New(config.StoragePath)
	s := server.New(config.Address, config.BaseURL, repo, db)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := s.Run(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}

}
