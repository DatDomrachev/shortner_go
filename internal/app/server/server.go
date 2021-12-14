package server

import (
	"context"
	"github.com/DatDomrachev/shortner_go/internal/app/handlers"
	"github.com/DatDomrachev/shortner_go/internal/app/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

type Server interface {
	configureRouter() *chi.Mux
}

type srv struct {
	address string
	baseURL string
	repo    repository.Repositorier
}

func New(address string, baseURL string, repo repository.Repositorier) *srv {
	server := &srv{
		address: address,
		baseURL: baseURL,
		repo:    repo,
	}

	return server
}

func (s *srv) Run(ctx context.Context) (err error) {
	router := s.ConfigureRouter()
	serv := &http.Server{
		Addr:    s.address,
		Handler: router,
	}

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("listener failed:+%v\n", err)
			cancel()
		}

	}()

	<-ctx.Done()

	log.Print("Server Started")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := serv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")

	return

}

func (s *srv) ConfigureRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{Id}", handlers.SimpleReadHandler(s.repo))
	router.Post("/", handlers.SimpleWriteHandler(s.repo, s.baseURL))
	router.Post("/api/shorten", handlers.SimpleJSONHandler(s.repo, s.baseURL))
	return router
}
