package store

import (
	"encoding/json"
	"fmt"
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

func TestKeys(t *testing.T) {
	db := setup(t)
	keys := []string{"test1", "test2", "test3"}
	expected := [][]byte{}
	for _, key := range keys {
		expected = append(expected, []byte(key))
		err := db.Set(key, "test")
		assert.Nil(t, err)
	}

	actual, err := db.Keys()
	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

	teardown(t)
}

func TestImport(t *testing.T) {
	db := setup(t)

	items := []kv{}
	for i := range 10 {
		items = append(items, kv{fmt.Sprint(i), fmt.Sprint(i)})
	}

	itemsJson, err := json.Marshal(items)
	assert.Nil(t, err)

	err = db.Import(itemsJson)
	assert.Nil(t, err)

	keys, err := db.Keys()
	assert.Nil(t, err)
	for i := range keys {
		assert.Equal(t, items[i].Key, string(keys[i]))
	}

	teardown(t)
}

func TestExport(t *testing.T) {
	db := setup(t)

	items := []kv{}
	for i := range 1 {
		items = append(items, kv{fmt.Sprint(i), fmt.Sprint(i)})
		err := db.Set(fmt.Sprint(i), fmt.Sprint(i))
		assert.Nil(t, err)
	}

	expected := `[{"Key":"0","Val":"\"0\""}]`

	actual, err := db.Export()
	assert.Nil(t, err)
	assert.Equal(t, expected, string(actual))

	teardown(t)
}
