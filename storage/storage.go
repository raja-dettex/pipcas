package storage

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

var DefaultPathTransformFunc = func(key string) (PathKey, error) {
	return PathKey{PathName: key}, nil
}

func CASTransformFunc(key string) (PathKey, error) {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	blockSize := 5
	pathsLen := len(hashStr) / blockSize
	paths := make([]string, pathsLen)
	for i := 0; i < pathsLen; i++ {
		from, to := i*blockSize, i*blockSize+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}, nil
}

type PathKey struct {
	PathName string
	FileName string
}

func (pathKey PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", pathKey.PathName, pathKey.FileName)
}
func (pathKey PathKey) FirstPathName() string {
	fmt.Println(pathKey.PathName)
	paths := strings.Split(pathKey.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

type PathTransformFunc func(string) (PathKey, error)

type StorageOpts struct {
	PathTransform PathTransformFunc
}

type Storage struct {
	opts StorageOpts
}

func NewStorage(opts StorageOpts) *Storage {
	return &Storage{
		opts: opts,
	}
}

func (s *Storage) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, f)
	if err != nil {
		return nil, err
	}
	return buff, nil
}

func (s *Storage) Delete(key string) error {
	pathKey, err := s.opts.PathTransform(key)
	if err != nil {
		return err
	}
	return os.RemoveAll(pathKey.FirstPathName())
}

func (s *Storage) readStream(key string) (io.ReadCloser, error) {
	pathKey, err := s.opts.PathTransform(key)
	if err != nil {
		return nil, err
	}
	return os.Open(pathKey.FullPath())
}

func (s *Storage) WriteStream(key string, r io.Reader) (string, error) {
	pathKey, err := s.opts.PathTransform(key)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return "", err
	}
	buff := new(bytes.Buffer)
	io.Copy(buff, r)
	fullPath := pathKey.FullPath()
	file, err := os.Create(fullPath)
	defer file.Close()
	if err != nil {
		return "", err
	}
	n, err := io.Copy(file, buff)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("written (%d) bytes to disk (%s)\n", n, fullPath), nil

}
