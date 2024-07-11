package store

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func setup(t *testing.T) *DB {
	err := os.RemoveAll("test.db")
	assert.Nil(t, err)
	db, err := New("test.db")
	assert.Nil(t, err)
	return db
}

func teardown(t *testing.T) {
	err := os.RemoveAll("test.db")
	assert.Nil(t, err)
}

func TestPut(t *testing.T) {
	db := setup(t)
	err := db.Set("test", "test")
	assert.Nil(t, err)
	teardown(t)
}

func TestGet(t *testing.T) {
	db := setup(t)
	value := "test"
	expected, err := json.Marshal(value)
	assert.Nil(t, err)
	err = db.Set("test", value)
	assert.Nil(t, err)
	actual, err := db.Get("test")
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	teardown(t)
}

func TestDelete(t *testing.T) {
	db := setup(t)
	expected := []byte{}
	err := db.Set("test", "test")
	assert.Nil(t, err)
	err = db.Delete("test")
	assert.Nil(t, err)
	actual, err := db.Get("test")
	assert.Equal(t, expected, actual)
	teardown(t)
}

func TestUnmarshal(t *testing.T) {
	db := setup(t)
	type testType struct {
		Title string
	}
	actual := testType{}
	expected := testType{
		Title: "test",
	}
	err := db.Set("test", expected)
	assert.Nil(t, err)
	err = db.Unmarshal("test", &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
	teardown(t)
}
