package database

import (
	"context"
	"log"
	"time"
	"database/sql"
    _ "github.com/jackc/pgx/v4"

) 

type DBWorker interface {
	Ping() (bool)
}

type DataBase struct {
	conn *sql.DB
}

func New(databaseURL string) (*DataBase, error) {
	if databaseURL == "" {
		return nil, nil
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	dataBase := &DataBase {
		conn: db,
	}
	return dataBase, nil
}

func (db *DataBase) Ping() (bool) {
	var bgCtx = context.Background()		
	ctx, cancel := context.WithTimeout(bgCtx, 1*time.Second)
    defer cancel()
    err:= db.conn.PingContext(ctx)
    if err != nil {
       log.Print(err.Error())
       return false;
    }

    return true;
}