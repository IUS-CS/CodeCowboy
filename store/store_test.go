package store

import (
	"os"
	"testing"
)

func TestPut(t *testing.T) {
	if os.Remove("test.db") != nil {
		panic("Cannot remove test.db")
	}
	db, err := New("test.db", "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	err = db.Set("test", "test")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
}
