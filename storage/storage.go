package storage

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Storer is the interface that defines the contract for our persistent storage.
// By using an interface, we can easily swap out the database implementation later.
type Storer interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close() error
}

// LevelDBStore is a LevelDB implementation of the Storer interface.
type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore creates and returns a new LevelDBStore instance.
func NewLevelDBStore(path string) (*LevelDBStore, error) {
	opts := &opt.Options{
		// You can configure LevelDB options here.
	}
	db, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return nil, err
	}
	return &LevelDBStore{db: db}, nil
}

// Put inserts a key-value pair into the database.
func (s *LevelDBStore) Put(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

// Get retrieves a value by its key from the database.
func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

// Has checks if a key exists in the database.
func (s *LevelDBStore) Has(key []byte) (bool, error) {
	return s.db.Has(key, nil)
}

// Delete removes a key-value pair from the database.
func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// Close closes the database connection.
func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
