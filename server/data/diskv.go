package data

import (
	"fmt"

	"github.com/peterbourgon/diskv"
)

// Database is disk persistent key value database
type Database struct {
	DB *diskv.Diskv
}

// Init initializes the db
func Init() Database {

	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	d := diskv.New(diskv.Options{
		BasePath:     "my-data-dir",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

	return Database{DB: d}
}

// Save stores data
func (d Database) Save(key string, data []byte) error {

	if err := d.DB.Write(key, data); err != nil {
		return fmt.Errorf("error writing data to the database %s", err)
	}
	return nil
}

// Load loads data
func (d Database) Load(key string) ([]byte, error) {
	value, err := d.DB.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error reading data from the database %s", err)
	}
	return value, nil
}

// Erase deletes data
func (d Database) Erase(key string) error {
	if err := d.DB.Erase(key); err != nil {
		return fmt.Errorf("error deleting data from the database %s", err)
	}
	return nil
}
