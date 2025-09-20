package storage

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

// Storer is an interface for a key-value store.
type Storer interface {
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	Has(key []byte) (bool, error)
	Delete(key []byte) error
	Close() error
}

// LevelDBStore is a Storer implementation using LevelDB.
type LevelDBStore struct {
	db *leveldb.DB
}

// NewLevelDBStore creates a new LevelDBStore instance.
func NewLevelDBStore(path string) (*LevelDBStore, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %w", err)
	}
	return &LevelDBStore{db: db}, nil
}

func (s *LevelDBStore) Put(key, value []byte) error {
	return s.db.Put(key, value, nil)
}

func (s *LevelDBStore) Get(key []byte) ([]byte, error) {
	return s.db.Get(key, nil)
}

func (s *LevelDBStore) Has(key []byte) (bool, error) {
	return s.db.Has(key, nil)
}

func (s *LevelDBStore) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

func (s *LevelDBStore) Close() error {
	return s.db.Close()
}
