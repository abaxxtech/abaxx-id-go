package dwn

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Blob represents a binary large object
type Blob []byte

// DataStore
// MemoryDatastore implements the DataStore interface using in-memory storage
type MemoryDatastore struct {
	data map[DataCid]Blob

	associated map[MessageCid]DataCid
}

func NewMemoryDatastore() DataStore {
	return &MemoryDatastore{
		data:       map[DataCid]Blob{},
		associated: map[MessageCid]DataCid{},
	}
}

func (*MemoryDatastore) Open() error {
	return nil
}

func (*MemoryDatastore) Close() error {
	return nil
}

func (m *MemoryDatastore) Put(tenant Tenant, messageCid MessageCid, dataCid DataCid,
	dataStream io.Reader) (resultCid DataCid, dataSize int64, err error) {
	return "", 0, nil
}

func (m *MemoryDatastore) Get(Tenant, MessageCid, DataCid) (dataCid DataCid,
	dataSize int64, dataStream io.Reader,
	err error) {
	return "", 0, nil, nil
}

func (m *MemoryDatastore) Associate(Tenant, MessageCid, DataCid) (dataCid DataCid,
	dataSize int64, err error) {
	return "", 0, nil
}

func (m *MemoryDatastore) Delete(Tenant, MessageCid, DataCid) (err error) {
	return nil

}

func (m *MemoryDatastore) Clear() (err error) {
	m.associated = make(map[MessageCid]DataCid)
	m.data = make(map[DataCid]Blob)

	return nil
}

// MessageStore
type Event struct {
	cid       MessageCid
	indexable map[string]IndexableValue
}

type MemoryEventLog struct {
	events []Event
}

func NewMemoryEventLog() EventLog {
	return &MemoryEventLog{[]Event{}}
}

func (*MemoryEventLog) Open() error {
	return nil
}

func (*MemoryEventLog) Close() error {
	return nil
}

func (*MemoryEventLog) Append(Tenant, MessageCid, IndexableKeyValues) (err error) {
	return nil
}

func (*MemoryEventLog) GetEvents(Tenant) ([]string, error) {
	return []string{}, nil
}

func (*MemoryEventLog) QueryEvents(Tenant, []Filter, EventLogCursor) ([]string, error) {
	return []string{}, nil
}

func (*MemoryEventLog) DeleteEventsByCid(Tenant, []MessageCid) error {
	return nil
}

// Test purposes
func (l *MemoryEventLog) Clear() error {
	l.events = make([]Event, 10)

	return nil
}

type GenericMessage struct {
	descriptor Descriptor
	data       []byte
}

// MessageStore
type Message struct{}

type MemoryMessageStore struct {
	messages map[string]struct {
		message   interface{}
		indexable IndexableKeyValues
	}
}

func NewMemoryMessageStore() MessageStore {
	return &MemoryMessageStore{
		messages: make(map[string]struct {
			message   interface{}
			indexable IndexableKeyValues
		}),
	}
}

func (m *MemoryMessageStore) Put(tenant Tenant, message interface{}, indexes IndexableKeyValues) (err error) {
	// For GenericMessage, use the DataCid from the descriptor as the key
	if genericMsg, ok := message.(*GenericMessage); ok {
		key := string(tenant) + ":" + string(genericMsg.descriptor.DataCid)
		m.messages[key] = struct {
			message   interface{}
			indexable IndexableKeyValues
		}{
			message:   genericMsg,
			indexable: indexes,
		}
		return nil
	}

	// For map[string]interface{}, try to use "id" field as key
	if msgMap, ok := message.(map[string]interface{}); ok {
		if id, ok := msgMap["id"].(string); ok {
			key := string(tenant) + ":" + id
			m.messages[key] = struct {
				message   interface{}
				indexable IndexableKeyValues
			}{
				message:   msgMap,
				indexable: indexes,
			}
			return nil
		}
	}

	// If we can't extract a proper key, generate a random one
	randomID := fmt.Sprintf("msg-%d", time.Now().UnixNano())
	key := string(tenant) + ":" + randomID

	m.messages[key] = struct {
		message   interface{}
		indexable IndexableKeyValues
	}{
		message:   message,
		indexable: indexes,
	}

	return nil
}

func (m *MemoryMessageStore) Get(tenant Tenant, messageCid MessageCid) (msg interface{}, err error) {
	key := string(tenant) + ":" + string(messageCid)
	if stored, ok := m.messages[key]; ok {
		return stored.message, nil
	}

	// If not found with direct key, try scanning all messages for matching DataCid
	for k, v := range m.messages {
		if strings.HasPrefix(k, string(tenant)+":") {
			if genericMsg, ok := v.message.(*GenericMessage); ok {
				if string(genericMsg.descriptor.DataCid) == string(messageCid) {
					return genericMsg, nil
				}
			}
		}
	}

	return nil, nil // Not found
}

func (*MemoryMessageStore) Query(tenant Tenant, filters []Filter,
	sort MessageSort,
	pagination Pagination) (err error) {
	return nil
}

func (*MemoryMessageStore) Delete(Tenant, MessageCid) (err error) {
	return nil
}

// Test purposes
func (m *MemoryMessageStore) Clear() (err error) {
	m.messages = make(map[string]struct {
		message   interface{}
		indexable IndexableKeyValues
	})

	return nil
}

func (*MemoryMessageStore) Open() error {
	return nil
}

func (*MemoryMessageStore) Close() error {
	return nil
}
