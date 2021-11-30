package handlers

import (
	"net/http"
	"io/ioutil"
    "github.com/DatDomrachev/shortner_go/internal/app/repository"
)



func SimpleReadHandler(repo repository.Repositorier) func(w http.ResponseWriter, r *http.Request) {
    return  func(w http.ResponseWriter, r *http.Request) {   
	    fullURL, err := repo.Load(r.URL.Path)
	   	
	   	if err != nil {
	        http.Error(w, "Not found", http.StatusBadRequest)
	    } 


	    if (fullURL != "") {
			http.Redirect(w,r, fullURL, http.StatusTemporaryRedirect)
	    } 
		

		http.Error(w, "Not found", http.StatusBadRequest)
	}
}	


func SimpleWriteHandler(repo repository.Repositorier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
		    http.Error(w, err.Error(), 500)
		    return
	    }

	  
	   result, err := repo.Store(string(data))
	   if err != nil {
		    http.Error(w, err.Error(), 500)
		    return
	    }

	   resp := "http://localhost:8080/" + result

	   w.Header().Set("content-type", "application/json")
	   w.WriteHeader(http.StatusCreated)
	   
	   w.Write([]byte(resp))
	}

}