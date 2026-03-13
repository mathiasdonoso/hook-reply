package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func NewConnection() (*Database, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".hook-replay")
	configPath := filepath.Join(configDir, "events.db")
	isNew := false

	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			isNew = true

			err = os.MkdirAll(configDir, 0755)
			if err != nil {
				return nil, err
			}

			file, err := os.Create(configPath)
			if err != nil {
				return nil, err
			}

			defer file.Close()
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", configPath)
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
