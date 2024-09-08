package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransfromFunc,
	}
	s := NewStore(opts)
	return s
}

func tearDown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpic"
	pathKey := CASPathTransfromFunc(key)
	expectedpathname := "1ff35/21faa/379e9/6fa2f/01352/2a356/cc934/c23d2"
	expectedoriginal := "1ff3521faa379e96fa2f013522a356cc934c23d2"
	if pathKey.Pathname != expectedpathname {
		t.Errorf("have %s want %s", pathKey.Pathname, expectedpathname)
	}

	if pathKey.Filename != expectedoriginal {
		t.Errorf("have %s want %s", pathKey.Filename, expectedoriginal)
	}
	assert.Equal(t, pathKey.Pathname, expectedpathname)
	assert.Equal(t, pathKey.Filename, expectedoriginal)
}

func TestStore(t *testing.T) {
	s := newStore()

	defer tearDown(t, s)

	count := 10

	for i := 0; i < count; i++ {

		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg data")
		if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
			fmt.Println(err)
			t.Error(err)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)

		fmt.Println(string(b))
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			fmt.Print(ok)
			t.Errorf("expected to NOT have key %s", key)
		}
	}
}

func TestRemove(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransfromFunc,
	}
	s := NewStore(opts)
	key := "OK"
	data := []byte("some jpg data")
	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		fmt.Println(err)
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
