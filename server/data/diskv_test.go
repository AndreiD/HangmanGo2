package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TesetInit(t *testing.T) {

	DB := Init()

	t.Run("Test Db Init Works", func(t *testing.T) {
		assert.NotNil(t, DB)
	})

}

func TestReadWrite(t *testing.T) {

	DB := Init()

	t.Run("Test Read / Write", func(t *testing.T) {
		DB.Save("mykey", []byte("payload"))

		savedValue, err := DB.Load("mykey")
		assert.Nil(t, err)

		assert.Equal(t, "payload", string(savedValue))
	})

}

func TestBadWriteRead(t *testing.T) {

	DB := Init()

	t.Run("Test Bad Write", func(t *testing.T) {
		err := DB.Save("", []byte("payload"))

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "error writing data to the database empty key")
		}

		_, err = DB.Load("")

		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "error reading data from the database file does not exist")
		}

	})

}

func TestErase(t *testing.T) {

	DB := Init()

	t.Run("Test Bad Write", func(t *testing.T) {
		DB.Save("mykey", []byte("payload"))

		err := DB.Erase("mykey")
		assert.Nil(t, err)

		err = DB.Erase("")
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "error deleting data from the database bad key")
		}

	})

}
