package dwn

import (
	"database/sql"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// SQLStoreConfig holds configuration for SQL-based storage
type SQLStoreConfig struct {
	DriverName     string // e.g., "postgres"
	DataSourceName string // e.g., "postgres://user:password@localhost:5432/dbname?sslmode=disable"
}

// SQLStore implements MessageStore interface using SQL database
type SQLStore struct {
	db     *sql.DB
	config SQLStoreConfig
}

// NewSQLStore creates a new SQLStore instance
func NewSQLStore(config SQLStoreConfig) (*SQLStore, error) {
	return &SQLStore{
		config: config,
	}, nil
}

// Open initializes the database connection and schema
func (s *SQLStore) Open() error {
	db, err := sql.Open(s.config.DriverName, s.config.DataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	s.db = db

	// Create the messages table
	schema := `
	CREATE TABLE IF NOT EXISTS dwn_messages (
		id SERIAL PRIMARY KEY,
		tenant TEXT NOT NULL,
		message_cid TEXT NOT NULL,
		message_data BYTEA NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(tenant, message_cid)
	);

	CREATE INDEX IF NOT EXISTS idx_dwn_messages_tenant ON dwn_messages(tenant);
	CREATE INDEX IF NOT EXISTS idx_dwn_messages_message_cid ON dwn_messages(message_cid);
	`

	_, err = s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (s *SQLStore) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Clear removes all messages (for testing)
func (s *SQLStore) Clear() error {
	_, err := s.db.Exec("DELETE FROM dwn_messages")
	return err
}

// Put stores a new message
func (s *SQLStore) Put(tenant Tenant, message interface{}, indexes IndexableKeyValues) error {
	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Extract message CID from the message
	msgMap, ok := message.(map[string]interface{})
	if !ok {
		return fmt.Errorf("message must be a map[string]interface{}")
	}

	messageCid, ok := msgMap["cid"].(string)
	if !ok {
		return fmt.Errorf("message must have a cid field")
	}

	// Encode the message
	messageBytes, err := cbor.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	// Insert the message
	query := `
		INSERT INTO dwn_messages (tenant, message_cid, message_data)
		VALUES ($1, $2, $3)
		ON CONFLICT (tenant, message_cid) DO UPDATE 
		SET message_data = EXCLUDED.message_data
	`
	_, err = tx.Exec(query, tenant, messageCid, messageBytes)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return tx.Commit()
}

// Get retrieves a message by its CID
func (s *SQLStore) Get(tenant Tenant, messageCid MessageCid) (interface{}, error) {
	var messageBytes []byte
	query := `SELECT message_data FROM dwn_messages WHERE tenant = $1 AND message_cid = $2`
	err := s.db.QueryRow(query, tenant, messageCid).Scan(&messageBytes)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	var message interface{}
	err = cbor.Unmarshal(messageBytes, &message)
	if err != nil {
		return nil, fmt.Errorf("failed to decode message: %w", err)
	}

	return message, nil
}

// Add these constants at the package level
const (
	SortDirectionAsc  = "ASC"
	SortDirectionDesc = "DESC"
)

// Add this helper function
func getSortDirectionString(direction SortDirection) string {
	if direction == Descending {
		return SortDirectionDesc
	}
	return SortDirectionAsc
}

// Update the Query method (replace the existing sort direction handling)
func (s *SQLStore) Query(tenant Tenant, filters []Filter, sort MessageSort, pagination Pagination) error {
	query := `SELECT message_data FROM dwn_messages WHERE tenant = $1`
	args := []interface{}{tenant}
	paramCount := 1

	// Add filters
	for _, filter := range filters {
		paramCount++
		query += fmt.Sprintf(" AND %s = $%d", filter.Property(), paramCount)
		args = append(args, filter.Value())
	}

	// Add sorting if provided
	if sort.Property != "" {
		query += fmt.Sprintf(" ORDER BY %s %s", sort.Property, getSortDirectionString(sort.Direction))
	}

	// Add pagination if provided
	if pagination != (Pagination{}) {
		if pagination.Limit > 0 {
			query += fmt.Sprintf(" LIMIT %d", pagination.Limit)
		}
		if pagination.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", pagination.Offset)
		}
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	return nil
}

// Delete removes a message
func (s *SQLStore) Delete(tenant Tenant, messageCid MessageCid) error {
	query := `DELETE FROM dwn_messages WHERE tenant = $1 AND message_cid = $2`
	result, err := s.db.Exec(query, tenant, messageCid)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("message not found")
	}

	return nil
}
