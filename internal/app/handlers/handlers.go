package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/DatDomrachev/shortner_go/internal/app/repository"
	"io/ioutil"
	"net/http"
	"strings"
)


func SimpleReadHandler(repo repository.Repositorier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fullURL, err := repo.Load(r.URL.Path)

		if err != nil {
			http.Error(w, "Not found", http.StatusBadRequest)
			return
		}

		
		http.Redirect(w, r, fullURL, http.StatusTemporaryRedirect)
	}
}

func SimpleWriteHandler(repo repository.Repositorier, baseURL string, userToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		result, err := repo.Store(string(data), userToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}


		if strings.Contains(result, "conflict:"){
			w.WriteHeader(http.StatusConflict)
			result = strings.ReplaceAll(result, "conflict:", "")

		} else {
			w.WriteHeader(http.StatusCreated)
		}
		

		resp := baseURL + "/" + result

		w.Header().Set("content-type", "application/json")
		

		w.Write([]byte(resp))
	}
}

func SimpleJSONHandler(repo repository.Repositorier, baseURL string, userToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var url repository.Item

		if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		result, err := repo.Store(url.FullURL, userToken)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if strings.Contains(result, "conflict:"){
			w.WriteHeader(http.StatusConflict)
			result = strings.ReplaceAll(result, "conflict:", "")

		} else {
			w.WriteHeader(http.StatusCreated)
		}

		newResult := repository.Result{ShortURL: baseURL + "/" + result}

		w.Header().Set("content-type", "application/json")
		

		buf := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(buf).Encode(newResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	}
}

func AllMyURLSHandler(repo repository.Repositorier, baseURL string, userToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		items := repo.GetByUser(userToken)

		if len(items) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}


		for i := range items {
			items[i].ShortURL = baseURL + "/" + items[i].ShortURL
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)

		buf := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(buf).Encode(items); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	
	}		
}


func PingDB(repo repository.Repositorier) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		result := repo.PingDB()

		if !result {
			http.Error(w, "No connection to DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}		
}

func BatchHandler(repo repository.Repositorier, baseURL string, userToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []repository.CorrelationItem 
		
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.Unmarshal(body, &items)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		
		shortens, err := repo.BatchAll(items, userToken)
		
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}	

		for i := range shortens {
			shortens[i].ShortURL = baseURL + "/" + shortens[i].ShortURL
		}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

		buf := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(buf).Encode(shortens); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	
	}		
}