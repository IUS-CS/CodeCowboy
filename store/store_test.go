package store

import "testing"

func TestPut(t *testing.T) {
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
