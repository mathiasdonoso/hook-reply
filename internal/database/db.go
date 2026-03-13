package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type database struct {
	db *sql.DB
}

func NewConnection(config *connectionConfig) (*database, error) {
	err := os.MkdirAll(config.Path, 0700)
	if err != nil {
		return nil, err
	}

	databasePath := filepath.Join(config.Path, config.DatabaseName)
	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	d := &database{db}
	if err = d.BuildSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return d, nil
}

func (d *database) Close() error {
	return d.db.Close()
}

func (d *database) GetDB() *sql.DB {
	return d.db
}

func (d *database) BuildSchema() error {
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
