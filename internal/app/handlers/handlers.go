package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/DatDomrachev/shortner_go/internal/app/repository"
	"github.com/DatDomrachev/shortner_go/internal/app/database"
	"io/ioutil"
	"net/http"
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

		resp := baseURL + "/" + result

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

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

		newResult := repository.Result{ShortURL: baseURL + "/" + result}

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)

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


func PingDB(db database.DBWorker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		result := db.Ping()

		if !result {
			http.Error(w, "No connection to DB", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)

	}		
}

