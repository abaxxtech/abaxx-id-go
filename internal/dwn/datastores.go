package dwn

import (
	"io"
)



type DataStore interface {
	Open() error
	Close() error

	// Put a data blog into the data store.  The DataCID is calculated from the
	// dataStream, and returned.  If it doesn't match dataCid input,
	// an error results, and nothing is written.
	Put(tenant Tenant, messageCid MessageCid, dataCid DataCid,
		dataStream io.Reader) (resultCid DataCid, dataSize int64, err error)

	Get(Tenant, MessageCid, DataCid) (dataCid DataCid,
		dataSize int64, dataStream io.Reader,
		err error)

	Associate(Tenant, MessageCid, DataCid) (dataCid DataCid,
		dataSize int64, err error)

	Delete(Tenant, MessageCid, DataCid) (err error)

	// Test purposes
	Clear() (err error)
}

type EventLogCursor string

type EventLog interface {
	Open() error
	Close() error

	Append(Tenant, MessageCid, IndexableKeyValues) (err error)

	GetEvents(Tenant) ([]string, error)

	QueryEvents(Tenant, []Filter, EventLogCursor) ([]string, error)

	DeleteEventsByCid(Tenant, []MessageCid) error

	// Test purposes
	Clear() error
}

// What are we storing in the MessageStore? What's in the DwnMessage?
type StoredMessage struct {
	Authorization map[string]interface{}
        Descriptor map[string]interface{}
}

type MessageStore interface {
	// The typescript versions of these functions take a
	// MessageStoreOptions which only contains an 'AbortSignal' --
	// we would possibly do this via channels later.

	// interface[} -> DwnMessage
	Put(Tenant,
		interface{},
		IndexableKeyValues) (err error)

	Get(Tenant, MessageCid) (msg interface{}, err error)

	Query(tenant Tenant, filters []Filter,
		sort MessageSort,
		pagination Pagination) (err error)

	Delete(Tenant, MessageCid) (err error)

	// Test purposes
	Clear() (err error)

	Open() error
	Close() error
}

// The data store can index items for each message.
// These are the types to support that.

// The storage calls take a map of string keys to
// values.
type IndexableKeyValues map[string]IndexableValue

// Indexable values are string | number | boolean, but Go doesn't
// have 'sum' types or unions.  This is the best alternative
// according to https://www.jerf.org/iri/post/2917/
type IndexableValue interface {
	isIndexableValue()
}

type F float64
type I int64
type B bool
type S string

func (f F) isIndexableValue() {}
func (i I) isIndexableValue() {}
func (b B) isIndexableValue() {}
func (s S) isIndexableValue() {}

// A filter applies a given FilterValue to a given property.
type Filter interface {
	// The property this filter is for
	Property() string
	Value() FilterValue
}

// A Filter compares either:
//   - equality (exactly matches an indexable value)
//   - one of a list of equality
//   - a range filter, which is an operator and a RangeValue.
//     The range value is a subset of indexable values: it's
//     numbers and strings only.
type FilterValue interface {
	isFilterValue()
}

type EqualFilter struct {
	EqualTo IndexableValue
}

func (o EqualFilter) isFilterValue() {}

type OneOfFilter struct {
	OneOf []EqualFilter
}

func (o OneOfFilter) isFilterValue() {}

// A Range Filter is one of:
// - GT (some value)
// - LT (some value)
// - GTE (some value)
// - LTE (some value)
//
// The value is a subset of indexable values, since you can't do a
// meaningful comparison on booleans.
type RangeFilter interface {
	isRangeFilter()
	isFilterValue()
	RangeValue() RangeValue
}

// Range Filter pieces below:
type GT struct {
	GT RangeValue
}
type GTE struct {
	GTE RangeValue
}
type LT struct {
	LT RangeValue
}
type LTE struct {
	LTE RangeValue
}

func (o GT) isRangeFilter()  {}
func (o GT) isFilterValue()  {}
func (o GTE) isRangeFilter() {}
func (o GTE) isFilterValue() {}
func (o LT) isRangeFilter()  {}
func (o LT) isFilterValue()  {}
func (o LTE) isRangeFilter() {}
func (o LTE) isFilterValue() {}

// A range value is:
// `string | number` numbers are either float64 or int64.
type RangeValue interface {
	isRangeValue()
}

func (s S) isRangeValue() {}
func (i I) isRangeValue() {}
func (f F) isRangeValue() {}

