package handlers

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
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

	   resp := repo.GetBaseURL() + result

	   w.Header().Set("content-type", "application/json")
	   w.WriteHeader(http.StatusCreated)
	   
	   w.Write([]byte(resp))
	}
}	


func SimpleJSONHandler(repo repository.Repositorier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {	

		var url repository.Item

		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			 http.Error(w, err.Error(), 500)
		    return
		}

		result, err := repo.Store(url.FullURL)
	    
	    if err != nil {
		     http.Error(w, err.Error(), 500)
		     return
	    }

	   

	    newResult := repository.Result{ShortURL: repo.GetBaseURL() + result}

	    w.Header().Set("content-type", "application/json")
	    w.WriteHeader(http.StatusCreated)

	    buf := bytes.NewBuffer([]byte{})
	    if err := json.NewEncoder(buf).Encode(newResult); err != nil {
			 http.Error(w, err.Error(), 500)
		    return
		}
	   
	   w.Write(buf.Bytes())
	}
}