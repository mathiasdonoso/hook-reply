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

type ConnectionConfig struct {
	DatabaseName string
	Path         string
}

func NewConnectionConfig() (*ConnectionConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configDir := filepath.Join(homeDir, ".hook-replay")

	return &ConnectionConfig{
		DatabaseName: "events.db",
		Path:         configDir,
	}, nil
}

func NewConnection(config *ConnectionConfig) (*Database, error) {
	databasePath := filepath.Join(config.Path, config.DatabaseName)

	err := os.MkdirAll(config.Path, 0700)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", databasePath)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	d := &Database{db}
	if err = d.BuildSchema(); err != nil {
		db.Close()
		return nil, err
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
