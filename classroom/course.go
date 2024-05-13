package classroom

import (
	"cso/codecowboy/store"
	"errors"
	"github.com/dgraph-io/badger/v3"
)

type Course struct {
	db          *store.DB
	Name        string
	Students    Students
	Assignments Assignments
}

func New(db *store.DB, name string) (*Course, error) {
	c := &Course{db: db, Name: name}
	err := c.Populate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func All(db *store.DB) ([]*Course, error) {
	keys, err := db.Keys()
	if err != nil {
		return nil, err
	}
	courses := []*Course{}
	for _, k := range keys {
		course, err := New(db, string(k))
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func (c *Course) Save() error {
	return c.db.Set(c.Name, c)
}

func (c *Course) Populate() error {
	err := c.db.Unmarshal(c.Name, c)
	if !errors.Is(err, badger.ErrKeyNotFound) && err != nil {
		return err
	}
	return nil
}
