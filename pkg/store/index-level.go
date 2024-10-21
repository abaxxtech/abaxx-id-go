package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type EqualFilter interface{}
type KeyValues map[string]interface{}
type RangeFilter map[string]interface{}
type Filter = map[string]interface{}

type IndexLevelConfig struct {
	Location string
}

type IndexedItem struct {
	ItemID  string    `json:"itemId"`
	Indexes KeyValues `json:"indexes"`
}

const (
	INDEX_SUBLEVEL_NAME = "index"
	DELIMITER           = "\x00"
	NEGATIVE_PREFIX     = "!"
	NEGATIVE_OFFSET     = int64(^uint64(0) >> 1)
	MAX_INT_STRING_LEN  = 19 // Length of string representation of max int64
)

type IndexLevelOptions struct {
	Context context.Context
}

type IndexLevel struct {
	db     *leveldb.DB
	config IndexLevelConfig
}

func NewIndexLevel(config IndexLevelConfig) (*IndexLevel, error) {
	if config.Location == "" {
		return nil, errors.New("location is required")
	}

	// Open LevelDB
	db, err := leveldb.OpenFile(config.Location, &opt.Options{
		Filter: filter.NewBloomFilter(10),
	})
	if err != nil {
		return nil, err
	}

	return &IndexLevel{
		db:     db,
		config: config,
	}, nil
}

func (il *IndexLevel) Close() error {
	return il.db.Close()
}

func (il *IndexLevel) Clear() error {
	// Not directly supported; requires deleting all keys
	// Alternatively, close and reopen the database
	err := il.db.Close()
	if err != nil {
		return err
	}
	db, err := leveldb.RecoverFile(il.config.Location, nil)
	if err != nil {
		return err
	}
	il.db = db
	return nil
}

func isEmptyObject(obj KeyValues) bool {
	return len(obj) == 0
}

func (il *IndexLevel) Put(tenant string, itemId string, indexes KeyValues, options *IndexLevelOptions) error {
	if isEmptyObject(indexes) {
		return errors.New("index must include at least one valid indexable property")
	}

	batch := new(leveldb.Batch)

	for indexName, indexValue := range indexes {
		key := keySegmentJoin(encodeValue(indexValue), itemId)
		item := IndexedItem{
			ItemID:  itemId,
			Indexes: indexes,
		}
		itemBytes, err := json.Marshal(item)
		if err != nil {
			return err
		}
		partitionKey := il.createIndexPartitionKey(tenant, indexName, key)
		batch.Put([]byte(partitionKey), itemBytes)
	}

	// Reverse lookup
	indexesBytes, err := json.Marshal(indexes)
	if err != nil {
		return err
	}
	reverseKey := il.createReverseLookupKey(tenant, itemId)
	batch.Put([]byte(reverseKey), indexesBytes)

	return il.db.Write(batch, nil)
}

func (il *IndexLevel) Delete(tenant string, itemId string, options *IndexLevelOptions) error {
	indexes, err := il.getIndexes(tenant, itemId)
	if err != nil {
		// Item not found
		return nil
	}

	batch := new(leveldb.Batch)

	// Delete reverse lookup
	reverseKey := il.createReverseLookupKey(tenant, itemId)
	batch.Delete([]byte(reverseKey))

	for indexName, indexValue := range indexes {
		key := keySegmentJoin(encodeValue(indexValue), itemId)
		partitionKey := il.createIndexPartitionKey(tenant, indexName, key)
		batch.Delete([]byte(partitionKey))
	}

	return il.db.Write(batch, nil)
}

func (il *IndexLevel) Query(tenant string, filters []Filter, queryOptions QueryOptions, options *IndexLevelOptions) ([]string, error) {
	// Implement query logic
	// For simplicity, we'll just return an empty slice
	return []string{}, nil
}

// Helper functions

func (il *IndexLevel) createIndexPartitionKey(tenant, indexName, key string) string {
	indexPartitionName := getIndexPartitionName(indexName)
	return fmt.Sprintf("%s%s%s%s%s", tenant, DELIMITER, indexPartitionName, DELIMITER, key)
}

func (il *IndexLevel) createReverseLookupKey(tenant, itemId string) string {
	return fmt.Sprintf("%s%s%s%s%s", tenant, DELIMITER, INDEX_SUBLEVEL_NAME, DELIMITER, itemId)
}

func getIndexPartitionName(indexName string) string {
	return fmt.Sprintf("__%s__", indexName)
}

func keySegmentJoin(values ...string) string {
	return strings.Join(values, DELIMITER)
}

func encodeNumberValue(value int64) string {
	prefix := ""
	offset := int64(0)
	if value < 0 {
		prefix = NEGATIVE_PREFIX
		offset = NEGATIVE_OFFSET
	}
	return fmt.Sprintf("%s%0*d", prefix, MAX_INT_STRING_LEN, value+offset)
}

func encodeValue(value interface{}) string {
	switch v := value.(type) {
	case int:
		return encodeNumberValue(int64(v))
	case int64:
		return encodeNumberValue(v)
	case float64:
		return encodeNumberValue(int64(v))
	case string:
		return fmt.Sprintf("%q", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		bytes, _ := json.Marshal(v)
		return string(bytes)
	}
}

func (il *IndexLevel) getIndexes(tenant, itemId string) (KeyValues, error) {
	reverseKey := il.createReverseLookupKey(tenant, itemId)
	data, err := il.db.Get([]byte(reverseKey), nil)
	if err != nil {
		return nil, err
	}
	var indexes KeyValues
	err = json.Unmarshal(data, &indexes)
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

// QueryOptions and other structs/enums used in the code

type QueryOptions struct {
	Limit         int
	Cursor        string
	SortProperty  string
	SortDirection string
}

const (
	SortDirectionAscending  = "ascending"
	SortDirectionDescending = "descending"
)
