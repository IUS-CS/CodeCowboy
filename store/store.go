package store

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	bolt "go.etcd.io/bbolt"
	"sync"
)

type ErrKeyNotFound struct {
	key string
	err error
}

func (e ErrKeyNotFound) Error() string {
	return fmt.Sprintf("key %s not found, %v", e.key, e.err)
}

type DB struct {
	path   string
	bucket []byte

	lock *sync.Mutex
}

func (db *DB) open(path string) (*bolt.DB, error) {
	db.lock.Lock()
	return bolt.Open(path, 0600, &bolt.Options{})
}

func (db *DB) close(kv *bolt.DB) error {
	err := kv.Close()
	db.lock.Unlock()
	return err
}

func New(path, bucket string) (*DB, error) {
	db := &DB{path, []byte(bucket), new(sync.Mutex)}
	kv, err := db.open(db.path)
	if err != nil {
		return nil, err
	}
	defer db.close(kv)
	if err = kv.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(db.bucket)
		return err
	}); err != nil {
		return nil, err
	}
	return db, nil
}

func (db *DB) Set(key string, input any) error {
	kv, err := db.open(db.path)
	if err != nil {
		return err
	}
	defer db.close(kv)
	value, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return kv.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		return b.Put([]byte(key), value)
	})
}

func (db *DB) Get(key string) ([]byte, error) {
	kv, err := db.open(db.path)
	if err != nil {
		return nil, err
	}
	defer db.close(kv)
	out := []byte{}
	err = kv.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		out = b.Get([]byte(key))
		return nil
	})
	return out, err
}

func (db *DB) Delete(key string) error {
	kv, err := db.open(db.path)
	if err != nil {
		return err
	}
	defer db.close(kv)
	return kv.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		err = b.Delete([]byte(key))
		return err
	})
}

func (db *DB) Keys() ([][]byte, error) {
	kv, err := db.open(db.path)
	if err != nil {
		return nil, err
	}
	defer db.close(kv)
	keys := [][]byte{}
	err = kv.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(db.bucket)
		return b.ForEach(func(k, v []byte) error {
			keys = append(keys, k)
			return nil
		})
	})
	return keys, err
}

func (db *DB) Unmarshal(key string, dest any) error {
	log.Debugf("unmarshaling %s from %v", key, string(db.bucket))
	value, err := db.Get(key)
	if err != nil {
		return err
	}
	if len(value) == 0 {
		err = ErrKeyNotFound{key, fmt.Errorf("key not found: %s", key)}
		log.Error(err)
		return err
	}
	return json.Unmarshal(value, dest)
}

func (db *DB) Export() ([]byte, error) {
	kv, err := db.open(db.path)
	if err != nil {
		return nil, err
	}
	defer db.close(kv)
	data := map[string]string{}
	err = kv.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		return b.ForEach(func(k, v []byte) error {
			data[string(k)] = string(v)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	out, err := json.Marshal(data)
	return out, err
}

func (db *DB) Import(data []byte) error {
	kv, err := db.open(db.path)
	if err != nil {
		return err
	}
	defer db.close(kv)
	input := map[string]string{}
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}
	return kv.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(db.bucket)
		if err != nil {
			return err
		}
		for k, v := range input {
			err = b.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}
		return nil
	})
}
