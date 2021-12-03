package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
    "github.com/DatDomrachev/shortner_go/internal/app/config"
    "context"
    "os"
	"os/signal"
	"log"

)



func main() {

	config, err := config.GetConfig()
	if err != nil {
		log.Printf("failed to configurate:+%v\n", err)
	}

	repo := repository.New(config.BaseURL)
	s:= server.New(config.Address, repo)

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
