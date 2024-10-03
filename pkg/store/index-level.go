package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

const INDEX_SUBLEVEL_NAME = "index"

type IndexLevelConfig struct {
	Location            string
	CreateLevelDatabase func(string) (*leveldb.DB, error)
}

type IndexLevel struct {
	db     *leveldb.DB
	config IndexLevelConfig
	mu     sync.RWMutex
}

type IndexedItem struct {
	ItemID  string    `json:"itemId"`
	Indexes KeyValues `json:"indexes"`
}

type QueryOptions struct {
	SortDirection SortDirection
	SortProperty  string
	Limit         int
	Cursor        string
}

type SortDirection int

const (
	Ascending SortDirection = iota
	Descending
)

func NewIndexLevel(config IndexLevelConfig) *IndexLevel {
	if config.CreateLevelDatabase == nil {
		config.CreateLevelDatabase = createLevelDatabase
	}

	db, err := config.CreateLevelDatabase(config.Location)
	if err != nil {
		return nil
	}
	return &IndexLevel{
		config: config,
		db:     db,
	}
}

func (il *IndexLevel) Open() error {
	il.mu.Lock()
	defer il.mu.Unlock()
	// The leveldb.DB type doesn't have an Open method, so we don't need to call it.
	// The database is already opened when it's created in NewIndexLevel.
	// We can just return nil to indicate success.
	return nil
}

func (il *IndexLevel) Close() error {
	il.mu.Lock()
	defer il.mu.Unlock()
	return il.db.Close()
}

func (il *IndexLevel) Clear() error {
	il.mu.Lock()
	defer il.mu.Unlock()
	// Iterate through all keys and delete them
	iter := il.db.NewIterator(nil, nil)
	defer iter.Release()

	batch := new(leveldb.Batch)
	for iter.Next() {
		batch.Delete(iter.Key())
	}

	if err := iter.Error(); err != nil {
		return fmt.Errorf("error iterating through keys: %w", err)
	}

	return il.db.Write(batch, nil)
}

func (il *IndexLevel) Put(tenant, itemID string, indexes KeyValues, ctx context.Context) error {
	il.mu.Lock()
	defer il.mu.Unlock()

	if len(indexes) == 0 {
		return fmt.Errorf("index must include at least one valid indexable property")
	}

	var indexOps []leveldb.Batch

	for indexName, indexValue := range indexes {
		key := keySegmentJoin(encodeValues(indexValue), itemID)
		item := IndexedItem{ItemID: itemID, Indexes: indexes}
		value, err := json.Marshal(item)
		if err != nil {
			return err
		}

		partitionOp, err := il.createOperationForIndexPartition(tenant, indexName, BatchOperation{
			Type:  "put",
			Key:   key,
			Value: value,
		})
		if err != nil {
			return err
		}
		indexOps = append(indexOps, partitionOp)
	}

	partitionOp, err := il.createOperationForIndexesLookupPartition(tenant, BatchOperation{
		Type:  "put",
		Key:   itemID,
		Value: []byte(encodeValues(indexes)),
	})
	if err != nil {
		return err
	}
	indexOps = append(indexOps, partitionOp)

	tenantPartition, err := il.db.Partition(tenant)
	if err != nil {
		return err
	}
	return tenantPartition.Batch(indexOps, ctx)
}

func (il *IndexLevel) Delete(tenant, itemID string, ctx context.Context) error {
	il.mu.Lock()
	defer il.mu.Unlock()

	indexes, err := il.getIndexes(tenant, itemID)
	if err != nil {
		return err
	}

	var indexOps []BatchOperation

	partitionOp, err := il.createOperationForIndexesLookupPartition(tenant, BatchOperation{
		Type: "del",
		Key:  itemID,
	})
	if err != nil {
		return err
	}
	indexOps = append(indexOps, partitionOp)

	for indexName, sortValue := range indexes {
		partitionOp, err := il.createOperationForIndexPartition(tenant, indexName, BatchOperation{
			Type: "del",
			Key:  keySegmentJoin(encodeValues(sortValue), itemID),
		})
		if err != nil {
			return err
		}
		indexOps = append(indexOps, partitionOp)
	}

	tenantPartition, err := il.db.Partition(tenant)
	if err != nil {
		return err
	}
	return tenantPartition.Batch(indexOps, ctx)
}

func (il *IndexLevel) Query(tenant string, filters []Filter, queryOptions QueryOptions, ctx context.Context) ([]string, error) {
	il.mu.RLock()
	defer il.mu.RUnlock()

	if shouldQueryWithInMemoryPaging(filters, queryOptions) {
		return il.queryWithInMemoryPaging(tenant, filters, queryOptions, ctx)
	}
	return il.queryWithIteratorPaging(tenant, filters, queryOptions, ctx)
}

// Add other methods like createOperationForIndexPartition, createOperationForIndexesLookupPartition,
// getIndexPartition, getIndexesLookupPartition, queryWithIteratorPaging, queryWithInMemoryPaging, etc.

func keySegmentJoin(values ...string) string {
	return strings.Join(values, "\x00")
}

func encodeValues(value interface{}) string {
	switch v := value.(type) {
	case float64:
		return encodeNumberValue(v)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func encodeNumberValue(value float64) string {
	const (
		NEGATIVE_OFFSET = float64(9007199254740991) // Number.MAX_SAFE_INTEGER
		NEGATIVE_PREFIX = "!"
		PADDING_LENGTH  = 16
	)

	prefix := ""
	offset := float64(0)
	if value < 0 {
		prefix = NEGATIVE_PREFIX
		offset = NEGATIVE_OFFSET
	}

	return fmt.Sprintf("%s%0*d", prefix, PADDING_LENGTH, int64(value+offset))
}
