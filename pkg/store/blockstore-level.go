package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/syndtr/goleveldb/leveldb"
)

// BlockstoreLevel implements the Blockstore interface using LevelDB
type BlockstoreLevel struct {
	db   *leveldb.DB
	path string
	mu   sync.RWMutex
}

// NewBlockstoreLevel creates a new BlockstoreLevel
func NewBlockstoreLevel(path string) (*BlockstoreLevel, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open leveldb: %w", err)
	}

	return &BlockstoreLevel{
		db:   db,
		path: path,
	}, nil
}

// Open opens the database
func (b *BlockstoreLevel) Open() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.db != nil {
		return nil
	}

	db, err := leveldb.OpenFile(b.path, nil)
	if err != nil {
		return fmt.Errorf("failed to open leveldb: %w", err)
	}

	b.db = db
	return nil
}

// Close closes the database
func (b *BlockstoreLevel) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.db == nil {
		return nil
	}

	err := b.db.Close()
	b.db = nil
	return err
}

// Put stores a block in the database
func (b *BlockstoreLevel) Put(ctx context.Context, c cid.Cid, block []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.db.Put(c.Bytes(), block, nil)
}

// Get retrieves a block from the database
func (b *BlockstoreLevel) Get(ctx context.Context, c cid.Cid) ([]byte, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.db.Get(c.Bytes(), nil)
}

// Has checks if a block exists in the database
func (b *BlockstoreLevel) Has(ctx context.Context, c cid.Cid) (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.db.Has(c.Bytes(), nil)
}

// Delete removes a block from the database
func (b *BlockstoreLevel) Delete(ctx context.Context, c cid.Cid) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.db.Delete(c.Bytes(), nil)
}

// PutMany stores multiple blocks in the database
func (b *BlockstoreLevel) PutMany(ctx context.Context, blocks map[cid.Cid][]byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	batch := new(leveldb.Batch)
	for c, block := range blocks {
		batch.Put(c.Bytes(), block)
	}
	return b.db.Write(batch, nil)
}

// AllKeysChan returns a channel of all CIDs in the database
func (b *BlockstoreLevel) AllKeysChan(ctx context.Context) (<-chan cid.Cid, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	ch := make(chan cid.Cid, 100)
	go func() {
		defer close(ch)
		iter := b.db.NewIterator(nil, nil)
		defer iter.Release()

		for iter.Next() {
			select {
			case <-ctx.Done():
				return
			default:
				c, err := cid.Cast(iter.Key())
				if err == nil {
					ch <- c
				}
			}
		}
	}()

	return ch, nil
}

// Clear deletes all entries in the database
func (b *BlockstoreLevel) Clear() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	iter := b.db.NewIterator(nil, nil)
	defer iter.Release()

	batch := new(leveldb.Batch)
	for iter.Next() {
		batch.Delete(iter.Key())
	}

	return b.db.Write(batch, nil)
}

// IsEmpty checks if the database is empty
func (b *BlockstoreLevel) IsEmpty() (bool, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	iter := b.db.NewIterator(nil, nil)
	defer iter.Release()

	return !iter.Next(), nil
}

// Add a Partition method to BlockstoreLevel
func (b *BlockstoreLevel) Partition(tenant string) (*BlockstoreLevel, error) {
	path := fmt.Sprintf("%s/%s", b.path, tenant)
	return NewBlockstoreLevel(path)
}
