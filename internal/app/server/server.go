package server

import(
	
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "context"
    "net/http"
    "github.com/DatDomrachev/shortner_go/internal/app/handlers"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
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

func (s *srv)Run(ctx context.Context) (err error) {
	router := s.ConfigureRouter()
	serv := &http.Server{
		Addr:    s.address,
		Handler: router,
	}
	
	
	go func() {
		if err := serv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
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

 
func (s *srv)ConfigureRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{Id}",handlers.SimpleReadHandler(s.repo)) 
	router.Post("/", handlers.SimpleWriteHandler(s.repo))
	return router 
}


