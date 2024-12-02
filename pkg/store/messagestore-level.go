package store

import (
	"context"

	"github.com/fxamacker/cbor/v2"
	"github.com/ipfs/go-cid"
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// MessageStoreLevel represents a message store using LevelDB
type MessageStoreLevel struct {
	blockstore *BlockstoreLevel
	index      *IndexLevel
	config     MessageStoreLevelConfig
}

// MessageStoreLevelConfig holds configuration for MessageStoreLevel
type MessageStoreLevelConfig struct {
	BlockstoreLocation  string
	IndexLocation       string
	CreateLevelDatabase func(string) (*LevelWrapper, error)
}

// MessageStoreOptions contains options for message store operations
type MessageStoreOptions struct {
	Signal context.Context
}

// SortDirection represents the direction of sorting
type SortDirection string

// Sort direction constants
const (
	SortAscending  SortDirection = "asc"
	SortDescending SortDirection = "desc"
)

// MessageSort specifies sorting options for messages
type MessageSort struct {
	DateCreated      *SortDirection
	DatePublished    *SortDirection
	MessageTimestamp *SortDirection
	Property         string // The property/column to sort by
	Direction        string // The sort direction ("ASC" or "DESC")

}

// Pagination specifies pagination options
type Pagination struct {
	Limit  int
	Cursor string
}

// NewMessageStoreLevel creates a new MessageStoreLevel instance
func NewMessageStoreLevel(config MessageStoreLevelConfig) (*MessageStoreLevel, error) {
	if config.CreateLevelDatabase == nil {
		config.CreateLevelDatabase = func(path string) (*LevelWrapper, error) {
			return createLevelDatabase(path), nil
		}
	}

	bs, err := NewBlockstoreLevel(config.BlockstoreLocation)
	if err != nil {
		return nil, err
	}

	idx, err := NewIndexLevel(IndexLevelConfig{
		Location: config.IndexLocation,
	})
	if err != nil {
		return nil, err
	}

	return &MessageStoreLevel{
		blockstore: bs,
		index:      idx,
		config:     config,
	}, nil
}

// Open opens the message store
func (msl *MessageStoreLevel) Open() error {
	return msl.blockstore.Open()
}

// Close closes the message store
func (msl *MessageStoreLevel) Close() error {
	return msl.blockstore.Close()
}

// Get retrieves a message by its CID
func (msl *MessageStoreLevel) Get(tenant, cidString string, options *MessageStoreOptions) (GenericMessage, error) {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return nil, options.Signal.Err()
		default:
		}
	}

	partition, err := msl.blockstore.Partition(tenant)
	if err != nil {
		return nil, err
	}

	c, err := cid.Decode(cidString)
	if err != nil {
		return nil, err
	}

	bytes, err := partition.Get(context.Background(), c)
	if err != nil {
		return nil, err
	}
	var message GenericMessage
	err = cbor.Unmarshal(bytes, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Query retrieves messages based on filters and sorting options
func (msl *MessageStoreLevel) Query(tenant string, filters []Filter, messageSort *MessageSort, pagination *Pagination, options *MessageStoreOptions) ([]GenericMessage, string, error) {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return nil, "", options.Signal.Err()
		default:
		}
	}

	queryOptions := buildQueryOptions(messageSort, pagination)
	results, err := msl.index.Query(tenant, filters, queryOptions, &IndexLevelOptions{})
	if err != nil {
		return nil, "", err
	}

	messages := make([]GenericMessage, 0, len(results))
	for _, messageCid := range results {
		message, err := msl.Get(tenant, messageCid, options)
		if err == nil {
			messages = append(messages, message)
		}
	}

	var cursor string
	if pagination != nil && pagination.Limit > 0 && len(results) > pagination.Limit {
		messages = messages[:pagination.Limit]
		cursor = results[pagination.Limit-1]
	}

	return messages, cursor, nil
}

// Delete removes a message from the store
func (msl *MessageStoreLevel) Delete(tenant, cidString string, options *MessageStoreOptions) error {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return options.Signal.Err()
		default:
		}
	}

	partition, err := msl.blockstore.Partition(tenant)
	if err != nil {
		return err
	}

	c, err := cid.Decode(cidString)
	if err != nil {
		return err
	}

	if err := partition.Delete(context.Background(), c); err != nil {
		return err
	}

	return msl.index.Delete(tenant, cidString, &IndexLevelOptions{})
}

// Put stores a new message in the store
func (msl *MessageStoreLevel) Put(tenant string, message GenericMessage, indexes KeyValues, options *MessageStoreOptions) error {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return options.Signal.Err()
		default:
		}
	}

	partition, err := msl.blockstore.Partition(tenant)
	if err != nil {
		return err
	}

	encodedMessage, err := cbornode.WrapObject(message, multihash.SHA2_256, -1)
	if err != nil {
		return err
	}

	messageCid := encodedMessage.Cid()
	if err := partition.Put(context.Background(), messageCid, encodedMessage.RawData()); err != nil {
		return err
	}

	messageCidString := messageCid.String()
	return msl.index.Put(tenant, messageCidString, indexes, &IndexLevelOptions{})
}

// Clear removes all messages from the store
func (msl *MessageStoreLevel) Clear() error {
	if err := msl.blockstore.Clear(); err != nil {
		return err
	}
	return msl.index.Clear()
}

func buildQueryOptions(messageSort *MessageSort, pagination *Pagination) QueryOptions {
	queryOptions := QueryOptions{
		SortDirection: string(SortAscending),
		SortProperty:  "messageTimestamp",
	}

	if messageSort != nil {
		if messageSort.DateCreated != nil {
			queryOptions.SortProperty = "dateCreated"
			queryOptions.SortDirection = string(*messageSort.DateCreated)
		} else if messageSort.DatePublished != nil {
			queryOptions.SortProperty = "datePublished"
			queryOptions.SortDirection = string(*messageSort.DatePublished)
		} else if messageSort.MessageTimestamp != nil {
			queryOptions.SortProperty = "messageTimestamp"
			queryOptions.SortDirection = string(*messageSort.MessageTimestamp)
		}
	}

	if pagination != nil {
		queryOptions.Limit = pagination.Limit
		if queryOptions.Limit > 0 {
			queryOptions.Limit++
		}
		queryOptions.Cursor = pagination.Cursor
	}

	return queryOptions
}

func createLevelDatabase(path string) *LevelWrapper {
	wrapper := NewLevelWrapper(LevelWrapperConfig{Location: path, OpenOptions: &opt.Options{NoSync: true}})
	return wrapper
}
