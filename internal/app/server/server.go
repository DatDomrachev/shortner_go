package server

import(
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "context"
    "net/http"
    "github.com/DatDomrachev/shortner_go/internal/app/handlers"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
    "os"
	"os/signal"
	"syscall"
	"log"
	"time"
)

type Server interface {
	configureRouter() *chi.Mux

}

type srv struct {
	address string
	repo repository.Repositorier
}

func NewServer(address string, repo repository.Repositorier) *srv{
	server := &srv {
		address: address,
		repo: repo, 
	} 
	
	return server
}

func (s *srv)Run() {
	router := s.ConfigureRouter()
	serv := &http.Server{
		Addr:    s.address,
		Handler: router,
	}
	
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	
	}()
	log.Print("Server Started")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

}

 
func (s *srv)ConfigureRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{Id}",handlers.SimpleReadHandler(s.repo)) 
	router.Post("/", handlers.SimpleWriteHandler(s.repo))
	return router 
}


