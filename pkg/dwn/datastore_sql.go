package dwn

import (
	"database/sql"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// SQLConfig holds configuration for SQL-based storage
type SQLConfig struct {
	DriverName     string
	DataSourceName string
}

// SQLStore implements MessageStore interface using SQL database
type SQLStore struct {
	db     *sql.DB
	config SQLConfig
}

// NewSQLStore creates a new SQLStore instance
func NewSQLStore(config SQLConfig) (*SQLStore, error) {
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

	// Create the message store table
	schema := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		tenant TEXT NOT NULL,
		message_cid TEXT NOT NULL,
		message_data BYTEA NOT NULL,
		interface TEXT,
		method TEXT,
		schema TEXT,
		data_cid TEXT,
		data_size TEXT,
		date_created TEXT,
		message_timestamp TEXT,
		data_format TEXT,
		is_latest_base_state TEXT,
		published TEXT,
		author TEXT,
		record_id TEXT,
		entry_id TEXT,
		date_published TEXT,
		latest TEXT,
		protocol TEXT,
		date_expires TEXT,
		description TEXT,
		granted_to TEXT,
		granted_by TEXT,
		granted_for TEXT,
		permissions_request_id TEXT,
		attester TEXT,
		protocol_path TEXT,
		recipient TEXT,
		context_id TEXT,
		parent_id TEXT,
		permissions_grant_id TEXT,
		UNIQUE(tenant, message_cid)
	)`

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

// Put stores a new message
func (s *SQLStore) Put(tenant Tenant, message interface{}, indexes IndexableKeyValues) error {
	// Encode message using CBOR
	messageBytes, err := cbor.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert message
	query := `
		INSERT INTO messages (
			tenant, message_cid, message_data,
			interface, method, schema, data_cid, data_size,
			date_created, message_timestamp, data_format,
			is_latest_base_state, published, author, record_id,
			entry_id, date_published, latest, protocol,
			date_expires, description, granted_to, granted_by,
			granted_for, permissions_request_id, attester,
			protocol_path, recipient, context_id, parent_id,
			permissions_grant_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31)
	`

	// Convert IndexableKeyValues to strings
	args := []interface{}{
		tenant,
		message.(map[string]interface{})["cid"],
		messageBytes,
		indexes["interface"],
		indexes["method"],
		indexes["schema"],
		indexes["dataCid"],
		indexes["dataSize"],
		indexes["dateCreated"],
		indexes["messageTimestamp"],
		indexes["dataFormat"],
		indexes["isLatestBaseState"],
		indexes["published"],
		indexes["author"],
		indexes["recordId"],
		indexes["entryId"],
		indexes["datePublished"],
		indexes["latest"],
		indexes["protocol"],
		indexes["dateExpires"],
		indexes["description"],
		indexes["grantedTo"],
		indexes["grantedBy"],
		indexes["grantedFor"],
		indexes["permissionsRequestId"],
		indexes["attester"],
		indexes["protocolPath"],
		indexes["recipient"],
		indexes["contextId"],
		indexes["parentId"],
		indexes["permissionsGrantId"],
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return tx.Commit()
}

// Get retrieves a message by its CID
func (s *SQLStore) Get(tenant Tenant, messageCid MessageCid) (interface{}, error) {
	var messageBytes []byte
	query := `SELECT message_data FROM messages WHERE tenant = $1 AND message_cid = $2`
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

// Query retrieves messages based on filters
func (s *SQLStore) Query(tenant Tenant, filters []Filter, sort MessageSort, pagination Pagination) error {
	query := `SELECT message_data FROM messages WHERE tenant = $1`
	args := []interface{}{tenant}
	paramCount := 1

	// Add filters
	for _, filter := range filters {
		paramCount++
		query += fmt.Sprintf(" AND %s = $%d", filter.Property(), paramCount)
		args = append(args, filter.Value())
	}

	if sort.Property != "" && sort.Direction != 0 {
		direction := "ASC"
		if sort.Direction == Descending {
			direction = "DESC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", sort.Property, direction)
	}

	if pagination.Cursor != "" {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pagination.Limit, pagination.Offset)
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
	query := `DELETE FROM messages WHERE tenant = $1 AND message_cid = $2`
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

// Clear removes all messages (for testing)
func (s *SQLStore) Clear() error {
	_, err := s.db.Exec("DELETE FROM messages")
	return err
}
