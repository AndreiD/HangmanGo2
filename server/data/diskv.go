package data

import (
	"fmt"

	"github.com/peterbourgon/diskv"
)

// DiskV is using https://github.com/peterbourgon/diskv
type DiskV struct {
	DB *diskv.Diskv
}

// Init initializes the db
func (d DiskV) Init() {

	// Simplest transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	d.DB = diskv.New(diskv.Options{
		BasePath:     "my-data-dir",
		Transform:    flatTransform,
		CacheSizeMax: 1024 * 1024,
	})

}

// Save stores data
func (d DiskV) Save(key string, data []byte) error {
	if err := d.DB.Write(key, data); err != nil {
		return fmt.Errorf("error writing data to the database %s", err)
	}
	return nil
}

// Load loads data
func (d DiskV) Load(key string) ([]byte, error) {
	value, err := d.DB.Read(key)
	if err != nil {
		return nil, fmt.Errorf("error reading data from the database %s", err)
	}
	return value, nil
}

// Erase deletes data
func (d DiskV) Erase(key string) error {
	if err := d.Erase(key); err != nil {
		return fmt.Errorf("error deleting data from the database %s", err)
	}
	return nil
}
