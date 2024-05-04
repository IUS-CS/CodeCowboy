package store

import (
	"encoding/json"
	"github.com/charmbracelet/charm/kv"
)

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

func (db *DB) Set(key string, input any) error {
	value, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return db.KV.Set([]byte(key), value)
}

func (db *DB) Get(key string) ([]byte, error) {
	return db.KV.Get([]byte(key))
}

func (db *DB) Unmarshal(key string, dest any) error {
	value, err := db.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(value, dest)
}
