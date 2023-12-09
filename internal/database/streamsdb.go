package database

import (
	"database/sql"

	"github.com/jackinthebox52/bytestream/internal/paths"
	_ "github.com/mattn/go-sqlite3"
)

var DB_CONNECTION *sql.DB

func InitializeConnection() error {
	db, err := InitializeSqliteConnection()
	if err != nil {
		return err
	}
	DB_CONNECTION = db
	return nil
}

func InitializeSqliteConnection() (*sql.DB, error) {
	if dpath, err := paths.CompileDatabasePath("streams"); err != nil {
		return nil, err
	} else {
		db, err := sql.Open("sqlite3", dpath)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err != nil {
			return nil, err
		}
		return db, nil
	}
}
