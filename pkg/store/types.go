package store

import (
	"database/sql"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/config"
)

type GenericMessage interface{}

type DataFilter struct {
	Property string
	Operator string
	Value    interface{}
}

// MessageStoreSQL represents a message store using SQL database
type MessageStoreSQL struct {
	db     *sql.DB
	config MessageStoreSQLConfig
}

// MessageStoreSQLConfig holds configuration for MessageStoreSQL
type MessageStoreSQLConfig struct {
	DBConfig config.DBConfig
}

// type KeyValues map[string]interface{}
