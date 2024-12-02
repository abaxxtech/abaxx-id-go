package store

import "database/sql"

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
	DriverName     string
	DataSourceName string
}

// type KeyValues map[string]interface{}
