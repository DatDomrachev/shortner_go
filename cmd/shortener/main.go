package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
)



func main() {
	
	repo := repository.New()
	s:= server.NewServer("localhost:8080", repo)
	s.Run()
	
}
