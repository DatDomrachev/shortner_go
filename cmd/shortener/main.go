package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
    "context"
    "os"
	"os/signal"
	"log"

)



func main() {
	
	repo := repository.New()
	s:= server.New("localhost:8080", repo)

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
