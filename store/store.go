package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jmoiron/sqlx"

	_ "github.com/mattn/go-sqlite3"
)

type ErrKeyNotFound struct {
	key string
	err error
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %s not found, %v", e.key, e.err)
}

type DB struct {
	*sqlx.DB

	path string
}

func (db *DB) open(path string) error {
	d, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return err
	}
	db.DB = d
	return nil
}

func (db *DB) close() error {
	return db.Close()
}

func New(path string) (*DB, error) {
	db := &DB{path: path}
	err := db.open(db.path)
	if err != nil {
		return nil, err
	}
	db.MustExec(`create table if not exists kv (key text, val text, 
		primary key (key))`)
	return db, nil
}

func (db *DB) Set(key string, input any) error {
	value, err := json.Marshal(input)
	if err != nil {
		return err
	}
	_, err = db.Exec(`insert into kv (key, val) values (?, ?)
		on conflict do update set val=?`,
		key, string(value), string(value))
	return err
}

func (db *DB) Get(key string) ([]byte, error) {
	out := []byte{}
	err := db.DB.Get(&out, `select val from kv where key=?`, key)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrKeyNotFound{key: key}
	}
	return out, err
}

func (db *DB) Delete(key string) error {
	_, err := db.Exec(`delete from kv where key=?`, key)
	return err
}

func (db *DB) Keys() ([][]byte, error) {
	keys := [][]byte{}
	err := db.Select(&keys, `select key from kv`)
	return keys, err
}

func (db *DB) Unmarshal(key string, dest any) error {
	log.Debugf("unmarshaling %s from %v", key)
	value, err := db.Get(key)
	if errors.As(err, &ErrKeyNotFound{}) {
		return nil
	}
	if err != nil {
		log.Debugf("error getting key %s: %v", key, err)
		return err
	}
	return json.Unmarshal(value, dest)
}

func (db *DB) Export() ([]byte, error) {
	data := []struct {
		Key string
		Val string
	}{}
	err := db.Select(&data, `select key, val from kv`)

	out, err := json.Marshal(data)
	return out, err
}

func (db *DB) Import(data []byte) error {
	input := []struct {
		Key string
		Val string
	}{}
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	for _, v := range input {
		_, err := db.Exec(`insert into kv (key, val) values (?, ?)`,
			v.Key, v.Val)
		if err != nil {
			return err
		}
	}
	return nil
}
