package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"log"
)

type Repositorier interface {
	Load(shortURL string, userToken string) (string, error)
	Store(url string, userToken string) (string, error)
	GetByUser(userToken string) ([]MyItem)
}

type Item struct {
	FullURL string `json:"url"`
	UserToken string `json:"user_token"`
}

type MyItem struct {
	ShortURL string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}


type Result struct {
	ShortURL string `json:"result"`
}

type Repo struct {
	StoragePath string
	items       []Item
}

func New(storagePath string) *Repo {
	var items []Item

	repo := &Repo{
		StoragePath: storagePath,
		items:       items,
	}

	if storagePath != "" {
		err := repo.readFromFile()

		if err != nil {
			log.Fatalf("failed to Load file:+%v", err)
		}
	}


	return repo
}


func (r *Repo) GetByUser(user string) ([]MyItem) {

	var myItems []MyItem

	for i := range r.items {
		if user == r.items[i].UserToken {
		  myItem := MyItem{
		 		ShortURL: strconv.Itoa(i+1),
		 		OriginalURL: r.items[i].FullURL,
			}
			myItems = append(myItems, myItem)
		}
	}

	return myItems
}

func (r *Repo) Load(shortURL string, userToken string) (string, error) {

	param := strings.TrimPrefix(shortURL, `/`)

	id, err := strconv.Atoi(param)

	if err != nil {
		return "", err
	}

	for i := range r.items {
		log.Print(r.items[i].UserToken)
		if i == id-1 && r.items[i].UserToken == userToken {
			return r.items[i].FullURL, nil
		}
	}
	return "", err
}

func (r *Repo) Store(url string, userToken string) (string, error) {
	newItem := Item{FullURL: url, UserToken: userToken}
	r.items = append(r.items, newItem)
	result := len(r.items)

	if r.StoragePath != "" {

		err := r.writeToFile(newItem)

		if err != nil {
			return "", err
		}

	}

	return strconv.Itoa(result), nil
}


func (r *Repo) readFromFile() (error) {
	file, err := os.OpenFile(r.StoragePath, os.O_RDONLY|os.O_CREATE, 0777)

	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
	
		data := scanner.Bytes()

		item := Item{}
		err := json.Unmarshal(data, &item)

		if err != nil {
			return err
		}

		r.items = append(r.items, item)

	}

	return nil
}

func (r *Repo) writeToFile(newItem Item) error {

	data, err := json.Marshal(&newItem)

	if err != nil {
		return err
	}

	file, err := os.OpenFile(r.StoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	if _, err := writer.Write(data); err != nil {
		return err
	}

	if err := writer.WriteByte('\n'); err != nil {
		return err
	}

	return writer.Flush()

}