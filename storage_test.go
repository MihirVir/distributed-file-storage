package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestpic"
	pathname := CASPathTransfromFunc(key)
	expectedpathname := "1ff35/21faa/379e9/6fa2f/01352/2a356/cc934/c23d2"
	if pathname != expectedpathname {
		t.Errorf("have %s want %s", pathname, expectedpathname)
	}

	assert.Equal(t, pathname, expectedpathname)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransfromFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some ok bytes"))
	if err := s.WriteStream("OK", data); err != nil {
		fmt.Println(err)
		t.Error(err)
	}
}
