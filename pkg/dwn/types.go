package dwn

// Types that are part of the public interface of the DWN.
type DID string
type Tenant string
type MessageCid string
type DataCid string

// A Raw Dwn message is just a parsed JSON placeholder.
type RawDwnMessage map[string]interface{}

type Status struct {
	Code   int
	Detail string
}

// Update the DwnConfig struct
type DwnConfig struct {
	DidResolver        *DidResolver
	TenantGate         TenantGate
	MessageStore       MessageStore
	DataStore          DataStore
	EventLog           EventLog
	BlockstoreLocation string
}

// Sort options
type SortDirection int

const (
	// these values are from query-types.ts
	Descending SortDirection = -1
	Ascending  SortDirection = 1
)

type MessageSort struct {
	DateCreated      SortDirection
	DatePublished    SortDirection
	MessageTimestamp SortDirection
}

type Pagination struct {
	Cursor string
	Limit  int
}
