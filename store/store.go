package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/charm/kv"
	"github.com/dgraph-io/badger/v3"
)

type ErrKeyNotFound struct {
	key string
	err error
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %s not found, %v", e.key, e.err)
}

type DB struct {
	name string
}

func New(name string) (*DB, error) {
	return &DB{name}, nil
}

func (db *DB) Set(key string, input any) error {
	kv, err := kv.OpenWithDefaults(db.name)
	if err != nil {
		return err
	}
	defer kv.Close()
	value, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return kv.Set([]byte(key), value)
}

func (db *DB) Get(key string) ([]byte, error) {
	kv, err := kv.OpenWithDefaults(db.name)
	if err != nil {
		return nil, err
	}
	defer kv.Close()
	return kv.Get([]byte(key))
}

func (db *DB) Delete(key string) error {
	kv, err := kv.OpenWithDefaults(db.name)
	if err != nil {
		return err
	}
	defer kv.Close()
	return kv.Delete([]byte(key))
}

func (db *DB) Keys() ([][]byte, error) {
	kv, err := kv.OpenWithDefaults(db.name)
	if err != nil {
		return nil, err
	}
	defer kv.Close()
	return kv.Keys()
}

func (db *DB) Unmarshal(key string, dest any) error {
	value, err := db.Get(key)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return ErrKeyNotFound{key, err}
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(value, dest)
}

func (db *DB) Export() ([]byte, error) {
	kv, err := kv.OpenWithDefaults(db.name)
	if err != nil {
		return nil, err
	}
	defer kv.Close()
	data := map[string]string{}
	keys, err := kv.Keys()
	if err != nil {
		return nil, err
	}
	for _, k := range keys {
		value, err := kv.Get(k)
		if err != nil {
			return nil, err
		}
		data[string(k)] = string(value)
	}
	out, err := json.Marshal(data)
	return out, err
}

func (db *DB) Import(data []byte) error {
	kv, err := kv.OpenWithDefaults(db.name)
	if err != nil {
		return err
	}
	defer kv.Close()
	input := map[string]string{}
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	for k, v := range input {
		err = kv.Set([]byte(k), []byte(v))
		if err != nil {
			return err
		}
	}
	return nil
}
