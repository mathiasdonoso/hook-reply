package database

import (
	"os"
	"path/filepath"
)

type connectionConfig struct {
	Opts
}

type OptFunc func(*Opts)

type Opts struct {
	DatabaseName string
	Path         string
}

func WithPath(path string) OptFunc {
	return func(opts *Opts) {
		opts.Path = path
	}
}

func WithDatabaseName(name string) OptFunc {
	return func(opts *Opts) {
		opts.DatabaseName = name
	}
}

func defaultOpts() (Opts, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return Opts{}, err
	}

	configDir := filepath.Join(homeDir, ".hook-replay")

	return Opts{
		DatabaseName: "events.db",
		Path:         configDir,
	}, nil
}

func NewConfig(opts ...OptFunc) (*connectionConfig, error) {
	o, err := defaultOpts()
	if err != nil {
		return nil, err
	}

	for _, fn := range opts {
		fn(&o)
	}

	return &connectionConfig{
		Opts: o,
	}, nil
}
