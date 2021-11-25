package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "context"
)



func main() {
	type contextKey string

	
	ctx := context.Background()
	s:= server.NewServer("localhost:8080")
	s.Run(ctx)
	
}
