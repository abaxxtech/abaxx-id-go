package store

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

type MessageStoreLevel struct {
	blockstore *BlockstoreLevel
	index      *IndexLevel
	config     MessageStoreLevelConfig
}

type MessageStoreLevelConfig struct {
	BlockstoreLocation  string
	IndexLocation       string
	CreateLevelDatabase func(string) (*LevelWrapper, error)
}

type MessageStoreOptions struct {
	Signal context.Context
}

type MessageSort struct {
	DateCreated      *SortDirection
	DatePublished    *SortDirection
	MessageTimestamp *SortDirection
}

type Pagination struct {
	Limit  int
	Cursor string
}

func NewMessageStoreLevel(config MessageStoreLevelConfig) *MessageStoreLevel {
	if config.CreateLevelDatabase == nil {
		config.CreateLevelDatabase = createLevelDatabase
	}

	return &MessageStoreLevel{
		blockstore: func() *BlockstoreLevel {
			bs, err := NewBlockstoreLevel(config.BlockstoreLocation)
			if err != nil {
				// Handle the error appropriately, e.g., log it or panic
				panic(err)
			}
			return bs
		}(),
		index: NewIndexLevel(IndexLevelConfig{
			Location:            config.IndexLocation,
			CreateLevelDatabase: config.CreateLevelDatabase,
		}),
		config: config,
	}
}

func (msl *MessageStoreLevel) Open() error {
	if err := msl.blockstore.Open(); err != nil {
		return err
	}
	return msl.index.Open()
}

func (msl *MessageStoreLevel) Close() error {
	if err := msl.blockstore.Close(); err != nil {
		return err
	}
	return msl.index.Close()
}

func (msl *MessageStoreLevel) Get(tenant string, cidString string, options *MessageStoreOptions) (GenericMessage, error) {
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

	bytes, err := partition.Get(c)
	if err != nil {
		return nil, err
	}

	var message GenericMessage
	err = cbornode.DecodeInto(bytes, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (msl *MessageStoreLevel) Query(tenant string, filters []Filter, messageSort *MessageSort, pagination *Pagination, options *MessageStoreOptions) ([]GenericMessage, string, error) {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return nil, "", options.Signal.Err()
		default:
		}
	}

	queryOptions := buildQueryOptions(messageSort, pagination)
	results, err := msl.index.Query(tenant, filters, queryOptions, options.Signal)
	if err != nil {
		return nil, "", err
	}

	var messages []GenericMessage
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

func (msl *MessageStoreLevel) Delete(tenant string, cidString string, options *MessageStoreOptions) error {
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

	if err := partition.Delete(c); err != nil {
		return err
	}

	return msl.index.Delete(tenant, cidString, options.Signal)
}

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
	if err := partition.Put(messageCid, encodedMessage.RawData()); err != nil {
		return err
	}

	messageCidString := messageCid.String()
	return msl.index.Put(tenant, messageCidString, indexes, options.Signal)
}

func (msl *MessageStoreLevel) Clear() error {
	if err := msl.blockstore.Clear(); err != nil {
		return err
	}
	return msl.index.Clear()
}

func buildQueryOptions(messageSort *MessageSort, pagination *Pagination) QueryOptions {
	queryOptions := QueryOptions{
		SortDirection: SortDirectionAscending,
		SortProperty:  "messageTimestamp",
	}

	if messageSort != nil {
		if messageSort.DateCreated != nil {
			queryOptions.SortProperty = "dateCreated"
			queryOptions.SortDirection = *messageSort.DateCreated
		} else if messageSort.DatePublished != nil {
			queryOptions.SortProperty = "datePublished"
			queryOptions.SortDirection = *messageSort.DatePublished
		} else if messageSort.MessageTimestamp != nil {
			queryOptions.SortProperty = "messageTimestamp"
			queryOptions.SortDirection = *messageSort.MessageTimestamp
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

func createLevelDatabase(path string) (*LevelWrapper, error) {
	return NewLevelWrapper(&LevelWrapperConfig{
		Location: path,
	})
}
