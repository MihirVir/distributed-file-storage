package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
)

const defaultRootFolderName = "mihir"

func CASPathTransfromFunc(key string) PathKey {
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

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashString,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FirstPathName() string {
	path := strings.Split(p.Pathname, "/")

	if len(path) == 0 {
		return ""
	}

	return path[0]
}

func (p PathKey) AbsolutePath() string {
	return fmt.Sprintf("%s%s", p.Pathname, p.Filename)
}

type StoreOpts struct {
	// folder name of the root, containing all the folders/files of the system
	Root              string
	PathTransformFunc PathTransformFunc
}

type Store struct {
	StoreOpts
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}

	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) WriteStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathnameWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.Pathname)

	if err := os.MkdirAll(pathnameWithRoot, os.ModePerm); err != nil {
		fmt.Println(reflect.TypeOf(err))
		return err
	}

	destination := pathKey.AbsolutePath()
	destinationWithRoot := fmt.Sprintf("%s/%s", s.Root, destination)
	f, err := os.Create(destinationWithRoot)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, destinationWithRoot)

	return nil
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

func (s *Store) Write(key string, r io.Reader) error {
	return s.WriteStream(key, r)
}

func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	pathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.AbsolutePath())
	return os.Open(pathKeyWithRoot)
}

func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("delete [%s] from disk", pathKey.Filename)
	}()

	return os.RemoveAll(s.Root + "/" + pathKey.FirstPathName())
}

func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, pathKey.AbsolutePath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}
