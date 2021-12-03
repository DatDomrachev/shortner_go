package repository

import (
	"strconv"
	"strings"
)

type Repositorier interface {
  Load(shortURL string) (string, error)
  Store(url string) (string, error)
  GetBaseURL() string
}


type Item struct {
	FullURL  string `json:"url"`
}

type Result struct {
	ShortURL string `json:"result"`
}

type Repo struct {
	BaseURL string
	items []Item
}


func New(baseURL string)*Repo{
	var items []Item

	repo := &Repo{
		BaseURL: baseURL,
		items: items,
		}
	return repo 
}

func (r *Repo) Load(shortURL string) (string, error) {
	param := strings.TrimPrefix(shortURL, `/`)

	id, err := strconv.Atoi(param)

	if err != nil {
        return "", err
    }  

	for i := range r.items {
		if i == id-1 {
			return r.items[i].FullURL, nil
		}	
	}
	return "", err	
}

func (r *Repo) Store(url string) (string, error) {
	newItem := Item{FullURL: url}
    r.items = append(r.items, newItem)	
    result := len(r.items)
   	return strconv.Itoa(result), nil
}

func (r *Repo) GetBaseURL() string {
	return r.BaseURL
}