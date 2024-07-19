package dwn

import (
	"io"
)

type Blob []byte

// DataStore
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

// MessageStore
type Message struct {
	cid     MessageCid
	indexes map[string]IndexableValue
}

type MemoryMessageStore struct {
	messages map[MessageCid]IndexableValue
}

func NewMemoryMessageStore() MessageStore {
	return &MemoryMessageStore{
		messages: map[MessageCid]IndexableValue{},
	}
}

func (*MemoryMessageStore) Put(Tenant,
	interface{},
	IndexableKeyValues) (err error) {
	return nil
}

func (*MemoryMessageStore) Get(Tenant, MessageCid) (msg interface{}, err error) {
	return nil, nil
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
	m.messages = make(map[MessageCid]IndexableValue)

	return nil
}

func (*MemoryMessageStore) Open() error {
	return nil
}

func (*MemoryMessageStore) Close() error {
	return nil
}
