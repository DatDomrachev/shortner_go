package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
    "context"
)



func main() {
	
	repo := repository.New()
	ctx := context.Background()
	s:= server.NewServer("localhost:8080", repo)
	s.Run(ctx)
	
}
