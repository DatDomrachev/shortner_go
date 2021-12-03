package handlers

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"bytes"
	"log"
	"encoding/json"
	"github.com/DatDomrachev/shortner_go/internal/app/repository"
)



func testRequest(t *testing.T, repo *repository.Repo, method, path, body string) (*http.Response, string) {
	
	request := httptest.NewRequest(method, "http://localhost:8080"+path, nil) 

	if (body != "") {
		var content = []byte(body)
		reqContent := bytes.NewBuffer(content)
		request = httptest.NewRequest(method, "http://localhost:8080"+path, reqContent)
	}
	

	//result, err := http.DefaultClient.Do(request)
	w := httptest.NewRecorder()
	
	
	if (method == "POST" && path =="/") {	
		SimpleWriteHandler(repo)(w, request)
	}

	if(method == "POST" && path =="/api/shorten") {
		SimpleJSONHandler(repo)(w, request)
	}

	if (method == "GET") {
		SimpleReadHandler(repo)(w, request)
	}	
	
	result := w.Result()

	respBody, err := ioutil.ReadAll(result.Body)
	require.NoError(t, err);

	defer result.Body.Close()
	return result, string(respBody)
	
}

func TestRouter(t *testing.T) {
	repo:=repository.New()

	result1, body1 := testRequest(t, repo, "POST", "/", "http://google.com")
	assert.Equal(t, 201, result1.StatusCode);
	assert.Equal(t, "application/json", result1.Header.Get("Content-Type"));	
	assert.Equal(t, "http://localhost:8080/1", body1);
	defer result1.Body.Close()	

	result2, body2 := testRequest(t, repo, "GET", "/1", "")
	assert.Equal(t, 307, result2.StatusCode);
	assert.Equal(t, "text/html; charset=utf-8", result2.Header.Get("Content-Type"));	
	assert.Equal(t, "http://google.com", result2.Header.Get("Location"));
	log.Println(body2)	
	defer result2.Body.Close()

	result3, body3 := testRequest(t, repo, "GET", "/aboba23","")
	assert.Equal(t, http.StatusBadRequest, result3.StatusCode);
	assert.Equal(t, "text/plain; charset=utf-8", result3.Header.Get("Content-Type"));	
	log.Println(body3)	
	defer result3.Body.Close()


	newQuery :=repository.Item{FullURL:"http://youtube.com"}	

	inputBuf := bytes.NewBuffer([]byte{})
    if err := json.NewEncoder(inputBuf).Encode(newQuery); err != nil {
		 log.Println(err.Error());
	    return
	}

	newResult := repository.Result{ShortURL:"http://localhost:8080/2"}	
	outputBuf := bytes.NewBuffer([]byte{})
    if err := json.NewEncoder(outputBuf).Encode(newResult); err != nil {
		log.Println(err.Error());
	    return
	}

	result4, body4 := testRequest(t, repo, "POST", "/api/shorten", inputBuf.String())
	assert.Equal(t, 201, result4.StatusCode);
	assert.Equal(t, "application/json", result4.Header.Get("Content-Type"));	
	assert.Equal(t, outputBuf.String(), body4);
	log.Println(body4)
	defer result4.Body.Close()	

}