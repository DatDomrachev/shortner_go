package repository

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"log"
	"context"
	"time"
	"database/sql"
  _ "github.com/jackc/pgx/v4/stdlib"
  "github.com/pressly/goose/v3"	 
)

type Repositorier interface {
	Load(shortURL string) (string, error)
	Store(url string, userToken string) (string, error)
	GetByUser(userToken string) ([]MyItem)
	PingDB()(bool)
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
	DatabaseURL string  
}

func New(storagePath string, databaseURL string) *Repo {
	var items []Item

	repo := &Repo{
		StoragePath: storagePath, 
		DatabaseURL: databaseURL,	
		items:       items,
	}

	if storagePath != "" {
		err := repo.readFromFile()

		if err != nil {
			log.Fatalf("failed to Load file:+%v", err)
		}
	}


	if databaseURL != "" {
		db, err := sql.Open("pgx", databaseURL)
		if err != nil {
			defer db.Close();
			log.Fatalf("Open DB Failed:%+v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("Open DB Failed:%+v", err)
		}

		

		err = goose.Up(db, "migrations" )
		if err != nil {
			log.Fatalf("failed executing migrations: %v\n", err)
		}
	}


	return repo
}


func (r *Repo) GetByUser(user string) ([]MyItem) {

	myItems:= make([]MyItem, 0)


	for i := range r.items {
		if user == r.items[i].UserToken {
		  myItem := MyItem{
		 		ShortURL: strconv.Itoa(i+1),
		 		OriginalURL: r.items[i].FullURL,
			}
			myItems = append(myItems, myItem)
		}
	}
	

	if r.DatabaseURL != "" {
		db, err := sql.Open("pgx", r.DatabaseURL)
		if err != nil {
			log.Print(err.Error())
			defer db.Close();
			return myItems
		}

		myItems = make([]MyItem, 0)
		ctx := context.Background()
		rows, err := db.QueryContext(ctx, "Select id::varchar(255), full_url from shortener.url WHERE user_token = $1", user)

		if err != nil {
			log.Print(err.Error())
			return myItems
		}

		defer rows.Close()

		for rows.Next() {
			var item MyItem
			err = rows.Scan(&item.ShortURL, &item.OriginalURL)

			if err != nil {
				log.Print(err.Error())
				return myItems
			}

			myItems = append(myItems, item)
		}

		err = rows.Err()
		if err != nil {
			log.Print(err.Error())
		}

	}

	return myItems
}

func (r *Repo) Load(shortURL string) (string, error) {

	fullURL := ""
	param := strings.TrimPrefix(shortURL, `/`)

	id, err := strconv.Atoi(param)


	if err != nil {
		return fullURL, err
	}

	for i := range r.items {
		if i == id-1 {
			fullURL = r.items[i].FullURL
			break
		}
	}

	if r.DatabaseURL != "" {
		db, err := sql.Open("pgx", r.DatabaseURL)
		if err != nil {
			log.Fatalf("Open DB Failed:%+v", err)
			defer db.Close();
			return "", err
		}
		err = db.QueryRow("SELECT full_url from shortener.url WHERE id = $1", id).Scan(&fullURL)
		if err != nil {
			log.Print(err.Error())
			return fullURL, err
		}
	}	

	return fullURL, nil
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

	if r.DatabaseURL != "" {
		db, err := sql.Open("pgx", r.DatabaseURL)
		if err != nil {
			log.Fatalf("Open DB Failed:%+v", err)
			defer db.Close();
			return "", err
		}

		err = db.QueryRow("Insert into shortener.url (full_url, user_token) VALUES ($1, $2) RETURNING id", url, userToken).Scan(&result)
		if err != nil {
			log.Fatalf("Insert DB Failed:%+v", err)
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

func (r *Repo) PingDB() (bool) {
	if r.DatabaseURL == "" {
		return false
	}

	db, err := sql.Open("pgx", r.DatabaseURL)
	if err != nil {
		log.Print(err.Error())
		defer db.Close();
		return false
	}
		
	
	var bgCtx = context.Background()		
	ctx, cancel := context.WithTimeout(bgCtx, 2*time.Second)
    defer cancel()
    err = db.PingContext(ctx)
    if err != nil {
       log.Print(err.Error())
       return false;
    }

    return true;
}