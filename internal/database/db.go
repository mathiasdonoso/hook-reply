package database

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewConnection() (*Database, error) {
	path := "./hr.db"
	isNew := false

	if _, err := os.Stat(path); err != nil {
		isNew = true
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	d := &Database{db}
	if isNew {
		err := d.BuildSchema()
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) BuildSchema() error {
	_, err := d.db.Exec(`
CREATE TABLE IF NOT EXISTS events (
	id STRING PRIMARY KEY,
	source TEXT NOT NULL,
	path TEXT NOT NULL,
	status TEXT NOT NULL,
	date TIMESTAMP NOT NULL
	);
	`)

	if err != nil {
		return err
	}

	return nil
}
