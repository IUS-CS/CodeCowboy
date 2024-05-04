package store

import "github.com/charmbracelet/charm/kv"

type DB struct {
	*kv.KV
}

func New(name string) (*DB, error) {
	db, err := kv.OpenWithDefaults(name)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
