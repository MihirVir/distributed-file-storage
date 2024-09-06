package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

func CASPathTransfromFunc(key string) string {
	hash := sha1.Sum([]byte(key))

	// hash[:] converts n number of bytes to slice
	hashString := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashString) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}

	return strings.Join(paths, "/")
}

type PathTransformFunc func(string) string

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathname := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathname, os.ModePerm); err != nil {
		fmt.Println(reflect.TypeOf(err))
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	filenameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(filenameBytes[:])
	destination := pathname + "/" + filename

	f, err := os.Create(destination)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, destination)

	return nil
}
