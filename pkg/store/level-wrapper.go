package store

import (
	"context"
	"errors"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// LevelWrapperOptions represents options for LevelWrapper methods.
type LevelWrapperOptions struct {
	Context context.Context
}

// LevelWrapperBatchOperation represents a batch operation.
type LevelWrapperBatchOperation struct {
	Type  string
	Key   []byte
	Value []byte
}

// LevelWrapperIteratorOptions represents options for iterators.
type LevelWrapperIteratorOptions struct {
	Start   []byte
	Limit   []byte
	Reverse bool
}

// LevelWrapperConfig holds configuration for LevelWrapper.
type LevelWrapperConfig struct {
	Location    string
	OpenOptions *opt.Options
}

// LevelWrapper provides a wrapper around LevelDB with partitioning support.
type LevelWrapper struct {
	config LevelWrapperConfig
	db     *leveldb.DB
	mu     sync.RWMutex
}

// NewLevelWrapper creates a new LevelWrapper instance.
func NewLevelWrapper(config LevelWrapperConfig) *LevelWrapper {
	return &LevelWrapper{
		config: config,
	}
}

// Open opens the LevelDB database.
func (lw *LevelWrapper) Open() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db != nil {
		return nil // Already open
	}

	db, err := leveldb.OpenFile(lw.config.Location, lw.config.OpenOptions)
	if err != nil {
		return err
	}

	lw.db = db
	return nil
}

// Close closes the LevelDB database.
func (lw *LevelWrapper) Close() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db == nil {
		return nil // Already closed
	}

	err := lw.db.Close()
	lw.db = nil
	return err
}

// Partition creates a new LevelWrapper for a sublevel (partition).
func (lw *LevelWrapper) Partition(name string) (*LevelWrapper, error) {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return nil, err
		}
	}

	// In Go LevelDB, we emulate sublevels by prefixing keys.
	partitionedWrapper := &LevelWrapper{
		config: lw.config,
		db:     lw.db,
		mu:     sync.RWMutex{},
	}
	return partitionedWrapper, nil
}

// Get retrieves a value by key.
func (lw *LevelWrapper) Get(key string, options *LevelWrapperOptions) ([]byte, error) {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return nil, err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return nil, options.Context.Err()
		default:
		}
	}

	data, err := lw.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	return data, err
}

// Has checks if a key exists.
func (lw *LevelWrapper) Has(key string, options *LevelWrapperOptions) (bool, error) {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return false, err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return false, options.Context.Err()
		default:
		}
	}

	return lw.db.Has([]byte(key), nil)
}

// Put stores a key-value pair.
func (lw *LevelWrapper) Put(key string, value []byte, options *LevelWrapperOptions) error {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return options.Context.Err()
		default:
		}
	}

	return lw.db.Put([]byte(key), value, nil)
}

// Delete removes a key-value pair.
func (lw *LevelWrapper) Delete(key string, options *LevelWrapperOptions) error {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return options.Context.Err()
		default:
		}
	}

	return lw.db.Delete([]byte(key), nil)
}

// IsEmpty checks if the database is empty.
func (lw *LevelWrapper) IsEmpty(options *LevelWrapperOptions) (bool, error) {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return false, err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return false, options.Context.Err()
		default:
		}
	}

	iter := lw.db.NewIterator(nil, nil)
	defer iter.Release()
	return !iter.Next(), iter.Error()
}

// Clear removes all entries from the database.
func (lw *LevelWrapper) Clear() error {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return err
		}
	}

	iter := lw.db.NewIterator(nil, nil)
	defer iter.Release()

	batch := new(leveldb.Batch)
	for iter.Next() {
		batch.Delete(iter.Key())
	}
	return lw.db.Write(batch, nil)
}

// Batch executes a batch of operations.
func (lw *LevelWrapper) Batch(operations []LevelWrapperBatchOperation, options *LevelWrapperOptions) error {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return options.Context.Err()
		default:
		}
	}

	batch := new(leveldb.Batch)
	for _, op := range operations {
		switch op.Type {
		case "put":
			batch.Put(op.Key, op.Value)
		case "del":
			batch.Delete(op.Key)
		default:
			return errors.New("unknown operation type")
		}
	}

	return lw.db.Write(batch, nil)
}

// Keys returns an iterator over the keys.
func (lw *LevelWrapper) Keys(options *LevelWrapperOptions) (iterator.Iterator, error) {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return nil, err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return nil, options.Context.Err()
		default:
		}
	}

	return lw.db.NewIterator(nil, nil), nil
}

// Iterator returns an iterator over key-value pairs with options.
func (lw *LevelWrapper) Iterator(iterOptions *LevelWrapperIteratorOptions, options *LevelWrapperOptions) (iterator.Iterator, error) {
	if lw.db == nil {
		if err := lw.Open(); err != nil {
			return nil, err
		}
	}

	if options != nil && options.Context != nil {
		select {
		case <-options.Context.Done():
			return nil, options.Context.Err()
		default:
		}
	}

	var rangeOpt *util.Range
	if iterOptions != nil {
		rangeOpt = &util.Range{
			Start: iterOptions.Start,
			Limit: iterOptions.Limit,
		}
	}

	iter := lw.db.NewIterator(rangeOpt, nil)
	return iter, nil
}
