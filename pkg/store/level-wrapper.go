package store

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelWrapperConfig struct {
	Location            string
	CreateLevelDatabase func(string) (*leveldb.DB, error)
	KeyEncoding         string
	ValueEncoding       string
}

type LevelWrapper struct {
	config LevelWrapperConfig
	db     *leveldb.DB
	mu     sync.RWMutex
}

type BatchOperation struct {
	Type  string
	Key   string
	Value []byte
}

func NewLevelWrapper(config *LevelWrapperConfig) (*LevelWrapper, error) {
	if config.CreateLevelDatabase == nil {
		config.CreateLevelDatabase = createLevelDatabase
	}
	db, err := config.CreateLevelDatabase(config.Location)
	if err != nil {
		return nil, err
	}
	return &LevelWrapper{
		config: *config,
		db:     db,
	}, nil
}

func (lw *LevelWrapper) Open() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db != nil {
		return nil
	}

	db, err := lw.config.CreateLevelDatabase(lw.config.Location)
	if err != nil {
		return fmt.Errorf("failed to open leveldb: %w", err)
	}

	lw.db = db
	return nil
}

func (lw *LevelWrapper) Close() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db == nil {
		return nil
	}

	err := lw.db.Close()
	lw.db = nil
	return err
}

func (lw *LevelWrapper) Partition(name string) (*LevelWrapper, error) {
	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.db == nil {
		return nil, fmt.Errorf("database not open")
	}

	return &LevelWrapper{
		config: LevelWrapperConfig{
			Location:            lw.config.Location + "/" + name,
			CreateLevelDatabase: lw.config.CreateLevelDatabase,
			KeyEncoding:         lw.config.KeyEncoding,
			ValueEncoding:       lw.config.ValueEncoding,
		},
		db: lw.db,
	}, nil
}

func (lw *LevelWrapper) Put(key string, value interface{}, ctx context.Context) error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db == nil {
		return fmt.Errorf("database not open")
	}

	encodedValue, err := encodeValue(value, lw.config.ValueEncoding)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return lw.db.Put([]byte(key), encodedValue, nil)
	}
}

func (lw *LevelWrapper) Get(key string, ctx context.Context) (interface{}, error) {
	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.db == nil {
		return nil, fmt.Errorf("database not open")
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		value, err := lw.db.Get([]byte(key), nil)
		if err == leveldb.ErrNotFound {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return decodeValue(value, lw.config.ValueEncoding)
	}
}

func (lw *LevelWrapper) Has(key string, ctx context.Context) (bool, error) {
	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.db == nil {
		return false, fmt.Errorf("database not open")
	}

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return lw.db.Has([]byte(key), nil)
	}
}

func (lw *LevelWrapper) Delete(key string, ctx context.Context) error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db == nil {
		return fmt.Errorf("database not open")
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return lw.db.Delete([]byte(key), nil)
	}
}

func (lw *LevelWrapper) Clear() error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db == nil {
		return fmt.Errorf("database not open")
	}

	iter := lw.db.NewIterator(nil, nil)
	defer iter.Release()

	batch := new(leveldb.Batch)
	for iter.Next() {
		batch.Delete(iter.Key())
	}

	return lw.db.Write(batch, nil)
}

func (lw *LevelWrapper) IsEmpty(ctx context.Context) (bool, error) {
	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.db == nil {
		return false, fmt.Errorf("database not open")
	}

	iter := lw.db.NewIterator(nil, nil)
	defer iter.Release()

	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
		return !iter.Next(), nil
	}
}

func (lw *LevelWrapper) Batch(operations []BatchOperation, ctx context.Context) error {
	lw.mu.Lock()
	defer lw.mu.Unlock()

	if lw.db == nil {
		return fmt.Errorf("database not open")
	}

	batch := new(leveldb.Batch)
	for _, op := range operations {
		switch op.Type {
		case "put":
			batch.Put([]byte(op.Key), op.Value)
		case "del":
			batch.Delete([]byte(op.Key))
		default:
			return fmt.Errorf("unknown operation type: %s", op.Type)
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return lw.db.Write(batch, nil)
	}
}

func (lw *LevelWrapper) Iterator(options *util.Range, ctx context.Context) *leveldb.Iterator {
	lw.mu.RLock()
	defer lw.mu.RUnlock()

	if lw.db == nil {
		return nil
	}

	return lw.db.NewIterator(options, nil)
}

func encodeValue(value interface{}, encoding string) ([]byte, error) {
	switch encoding {
	case "json":
		return json.Marshal(value)
	case "utf8":
		if s, ok := value.(string); ok {
			return []byte(s), nil
		}
		return nil, fmt.Errorf("value is not a string for utf8 encoding")
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func decodeValue(value []byte, encoding string) (interface{}, error) {
	switch encoding {
	case "json":
		var v interface{}
		err := json.Unmarshal(value, &v)
		return v, err
	case "utf8":
		return string(value), nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}
