package repository

import (
	"strconv"
	"strings"
)

type Repositorer interface {
  Load(shortURL string) (string, error)
  Store(url string) (string, error)
}


type item struct {
	FullURL  string
}

var items []item


func Load(shortURL string) (string, error) {
	param := strings.TrimPrefix(shortURL, `/`)

	id, err := strconv.Atoi(param)

	if err != nil {
        return "", err
    }  

	for i := range items {
		if i == id-1 {
			return items[i].FullURL, nil
		}	
	}

	return "", err	
}

func Store(url string) (string, error){
	newItem := item{FullURL: url}
    items = append(items, newItem)	
    result := len(items)
   	return strconv.Itoa(result), nil
}