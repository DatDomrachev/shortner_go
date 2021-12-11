package repository

import (
	"strconv"
	"strings"
	"os"
	"bufio"
	"encoding/json"
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
	StoragePath string
	items []Item
}


func New(baseURL string, storagePath string)*Repo{
	var items []Item

	repo := &Repo{
		BaseURL: baseURL,
		StoragePath: storagePath,
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


	if (r.StoragePath != "") {
	  result, err := r.readFromFile(id);
	
		if err != nil {
    	return "", err
  	}

  	return result, nil  
	}

	
	for i := range r.items {
		if i == id-1 {
			return r.items[i].FullURL, nil
		}	
	}
	return "", err	
}

func (r *Repo) Store(url string) (string, error) {
	  result := 0
	  newItem := Item{FullURL: url}
	 
	  if r.StoragePath == "" {
    	r.items = append(r.items, newItem)	
    	result = len(r.items)
  	} else {
   		
   		err := r.writeToFile(newItem)   
  		
  		if err != nil {
    		return "", err
  		}

  		result, err = r.countFileLines()

  		if err != nil {
    		return "", err
  		}
  	
  	}

   	
   	return strconv.Itoa(result), nil
}

func (r *Repo) GetBaseURL() string {
	return r.BaseURL
}


func (r* Repo) readFromFile(ID int) (string, error) {
	file, err := os.OpenFile(r.StoragePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
    return "", err
  }

  defer file.Close()

	scanner:= bufio.NewScanner(file)
	
	row :=0
  for scanner.Scan() {

	  if row == ID-1 {
	  	data := scanner.Bytes()

	  	item := Item{}
	  	err := json.Unmarshal(data, &item)
	  	
	  	if err != nil {
	    	  return "", err
	  	}

	  	return item.FullURL, nil
	  }

	  row++
	}

	return "", err	
}


func (r* Repo) writeToFile(newItem Item) error {
	
	data, err := json.Marshal(&newItem)
  
  if err != nil {
      return err
  }
  

  file, err := os.OpenFile(r.StoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	
	if err != nil {
        return err
  }

  defer file.Close()

   writer:= bufio.NewWriter(file)


  if _, err := writer.Write(data); err != nil {
        return err
  }

  if err := writer.WriteByte('\n'); err != nil {
      return err
  }

  return writer.Flush()

}


func (r* Repo) countFileLines() (int, error) {
	file, err := os.OpenFile(r.StoragePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
    return 0, err
  }
  
  defer file.Close()

  scanner:= bufio.NewScanner(file)

	row := 0
  for scanner.Scan() {
  	row++
  }

  return row, nil;

}