package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
    "context"
    "time"
    "os"
	"os/signal"
	"log"

)



func main() {
	
	repo := repository.New()
	s:= server.NewServer("localhost:8080", repo)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := s.Run(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}
	
}
