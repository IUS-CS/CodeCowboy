package classroom

import (
	"cso/codecowboy/store"
	"errors"
	"fmt"
)

type Course struct {
	db          *store.DB
	Name        string
	Students    Students
	Assignments Assignments
	Instructors []string
}

func New(db *store.DB, name string) (*Course, error) {
	c := &Course{db: db, Name: name}
	err := c.Populate()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Course) CanAccess(username string) bool {
	for _, instructor := range c.Instructors {
		if instructor == username {
			return true
		}
	}
	return false
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

func AllPermitted(db *store.DB, username string) ([]*Course, error) {
	courses, err := All(db)
	if err != nil {
		return nil, err
	}

	userCourses := []*Course{}
	for _, course := range courses {
		if course.CanAccess(username) {
			userCourses = append(userCourses, course)
		}
	}

	return userCourses, nil
}

func (c *Course) Validate() error {
	errs := []error{}
	if c.Name == "" {
		errs = append(errs, fmt.Errorf("course name is required"))
	}
	for _, a := range c.Assignments {
		if err := a.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("assignment %s validation failed: %w", a.Name, err))
		}
	}
	for _, s := range c.Students {
		if err := s.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("student %s is invalid: %w", s.Name, err))
		}
	}
	return errors.Join(errs...)
}

func (c *Course) Save() error {
	if err := c.Validate(); err != nil {
		return err
	}
	return c.db.Set(c.Name, c)
}

func (c *Course) Populate() error {
	err := c.db.Unmarshal(c.Name, c)
	notFoundError := store.ErrKeyNotFound{}
	if !errors.As(err, &notFoundError) && err != nil {
		return err
	}
	return nil
}
