package main

import (
    "github.com/DatDomrachev/shortner_go/internal/app/server"
    "context"
)



func main() {
	ctx := context.WithValue(context.Background(), "address", "localhost:8080")
	s:= server.NewServer("address")
	s.Run(ctx)
	
}
