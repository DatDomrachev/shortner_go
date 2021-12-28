package server

import (
	"context"
	"github.com/DatDomrachev/shortner_go/internal/app/handlers"
	"github.com/DatDomrachev/shortner_go/internal/app/repository"
	"github.com/DatDomrachev/shortner_go/internal/app/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
	"compress/gzip"
	"strings"
	"io"
	"io/ioutil"
	"bytes"
	"encoding/hex"
	"crypto/hmac"
    "crypto/sha256"
    "math/rand"
    "encoding/binary"
)

type Server interface {
	configureRouter() *chi.Mux
}

type srv struct {
	address string
	baseURL string
	repo    repository.Repositorier
	db    	database.DBWorker
}

type gzipWriter struct {
    http.ResponseWriter
    Writer io.Writer
}

type contextKey string

func (w gzipWriter) Write(b []byte) (int, error) {
    return w.Writer.Write(b)
} 

func New(address string, baseURL string, repo repository.Repositorier, db database.DBWorker) *srv {
	server := &srv{
		address: address,
		baseURL: baseURL,
		repo:    repo,
		db: 	 db,
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
	router.Use(GzipHandle)
	router.Use(CookieManager)

	router.Get("/{Id}", handlers.SimpleReadHandler(s.repo)) 
	router.Post("/", func(rw http.ResponseWriter, r *http.Request) {
    	u := r.Context().Value(contextKey("user_token")).(string)
    	handlers.SimpleWriteHandler(s.repo, s.baseURL, u)(rw,r)
    })
	router.Post("/api/shorten", func(rw http.ResponseWriter, r *http.Request) {
    	u := r.Context().Value(contextKey("user_token")).(string)
    	handlers.SimpleJSONHandler(s.repo, s.baseURL, u)(rw,r)
    })
	router.Get("/user/urls", func(rw http.ResponseWriter, r *http.Request) {
    	u := r.Context().Value(contextKey("user_token")).(string)
    	handlers.AllMyURLSHandler(s.repo, s.baseURL,u)(rw,r)
    })
    router.Get("/ping", handlers.PingDB(s.db))
	return router
}


func GzipHandle(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      
        if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
            next.ServeHTTP(w, r)
            return
        }

        

        gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
        if err != nil {
            io.WriteString(w, err.Error())
            return
        }
        defer gz.Close()

        if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
            reader, err := gzip.NewReader(r.Body)
       		
	        if err != nil {
	            io.WriteString(w, err.Error())
	            return
	        }

	        defer reader.Close()

	        b, err := ioutil.ReadAll(reader)

	        if err != nil {
	            io.WriteString(w, err.Error())
	            return
	        }


	        r.Body = ioutil.NopCloser(bytes.NewBuffer(b))    
        }

           

        w.Header().Set("Content-Encoding", "gzip")
        next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
    })
}


func newCookie(key []byte) (cookie *http.Cookie, err error) {
	
	//рандомный id
	src := make([]byte, 4)
	id:=rand.Uint32()
    binary.BigEndian.PutUint32(src, id)
	 
    
    //подпись	
    h := hmac.New(sha256.New, key)
    h.Write(src)

    cookie = &http.Cookie {
	        	Name:   "user_token",
	        	Value:  hex.EncodeToString(src) + hex.EncodeToString(h.Sum(nil)),
	        	MaxAge: 300,
	        }	

    return cookie, nil;
}

func CookieManager(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    	var secretkey = []byte("secret key of My Castle")


    	cookie, err := r.Cookie("user_token")

    	if (err == nil) {
    		
	    	userKey := cookie.Value

	    	data, err := hex.DecodeString(userKey)

	    	if err != nil {
		        log.Fatalf("CookieManager error:%+v", err)
		    }
		    
		   
		   
		    h := hmac.New(sha256.New, secretkey)
		    h.Write(data[:4])
		    sign := h.Sum(nil) 

		    if !hmac.Equal(sign, data[4:]) {
		       
		    	cookie, err = newCookie(secretkey)
		    	if err != nil {
				    log.Fatalf("CookieManager error:%+v", err)
				}
		    }

    	} else {

    		cookie, err = newCookie(secretkey)
    		if err != nil {
			    log.Fatalf("CookieManager error:%+v", err)
			}
    	}

    	

    	http.SetCookie(w, cookie);

    	ctx := context.WithValue(r.Context(), contextKey("user_token"), cookie.Value)
    	
    	next.ServeHTTP(w, r.WithContext(ctx))
    })
}    
