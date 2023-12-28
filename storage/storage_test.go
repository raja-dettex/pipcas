package storage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileStore(t *testing.T) {
	pathKey := "picture"
	storageOpts := StorageOpts{
		PathTransform: CASTransformFunc,
	}
	storage := NewStorage(storageOpts)
	data := []byte("some beautiful jpeg bytes ")
	res, err := storage.WriteStream(pathKey, bytes.NewBuffer(data))
	assert.Nil(t, err)
	fmt.Println(res)
	r, err := storage.Read(pathKey)
	assert.Nil(t, err)
	fetchedData, err := ioutil.ReadAll(r)
	assert.Equal(t, data, fetchedData)
	//err = storage.Delete(pathKey)
	fmt.Println(string(fetchedData))
	//fmt.Println(err.Error())
	//assert.Nil(t, err)
}
