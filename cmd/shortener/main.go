package main

import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
    "io/ioutil"
    "strconv"
    "log"
    "strings"
)

type URL struct {
	URL string `json:"url"`
}


type item struct {
	FullURL  string
}


var items []item


func SimpleReadHandler(w http.ResponseWriter, r *http.Request) {
	//param := chi.URLParam(r, "Id")
	param := strings.TrimPrefix(r.URL.Path, `/`)

	id, err := strconv.Atoi(param)

	if err != nil {
    	http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }  

   
	for i := range items {
		if i == id-1 {
			http.Redirect(w,r, items[i].FullURL, http.StatusTemporaryRedirect)
			return
		}	
	}	

	http.Error(w, "Not found", http.StatusBadRequest)
}


func SimpleWriteHandler(w http.ResponseWriter, r *http.Request) {	
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
	    http.Error(w, err.Error(), 500)
	    return
    }

   full := string(data)
   newItem := item{FullURL: full}
   items = append(items, newItem)	
   w.Header().Set("content-type", "application/json")
   w.WriteHeader(http.StatusCreated)
   result := len(items)
   resp := "http://localhost:8080/"+ strconv.Itoa(result)
   w.Write([]byte(resp))

}

func SimpleRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/{Id}",SimpleReadHandler) 
	router.Post("/", SimpleWriteHandler)

	return router
}

func main() {
	r := SimpleRouter()
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
