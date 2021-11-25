package server

import(
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "context"
    "net/http"
    "github.com/DatDomrachev/shortner_go/internal/app/handlers"
)

type Server interface {
	configureRouter() *chi.Mux

}

type srv struct {
	address string
}

func NewServer(address string) *srv{
	server := &srv {
		address: address,
	} 
	
	return server
}

func (s *srv)Run(ctx context.Context) error {
	router := ConfigureRouter()
	
	return http.ListenAndServe(s.address, router)
}

 
func ConfigureRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{Id}",handlers.SimpleReadHandler) 
	router.Post("/", handlers.SimpleWriteHandler)
	return router 
}


