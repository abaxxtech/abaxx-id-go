package store

import (
	"database/sql"
	"fmt"

	"github.com/fxamacker/cbor/v2"
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
)

// Filter represents a query filter
// type Filter struct {
// 	Property string      // The property/column to filter on
// 	Operator string      // The operator to use (e.g., "=", ">", "<", "LIKE")
// 	Value    interface{} // The value to filter by
// }

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

// NewMessageStoreSQL creates a new MessageStoreSQL instance
func NewMessageStoreSQL(config MessageStoreSQLConfig) (*MessageStoreSQL, error) {
	return &MessageStoreSQL{
		config: config,
	}, nil
}

// Open opens the message store and initializes the database schema
func (mss *MessageStoreSQL) Open() error {
	db, err := sql.Open(mss.config.DriverName, mss.config.DataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	mss.db = db

	// Create the message store table
	schema := `
	CREATE TABLE IF NOT EXISTS message_store (
		id SERIAL PRIMARY KEY,
		tenant TEXT NOT NULL,
		message_cid VARCHAR(60) NOT NULL,
		encoded_message_bytes BYTEA NOT NULL,
		encoded_data TEXT,
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
		permissions_grant_id TEXT
	)`

	_, err = mss.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}

	return nil
}

// Close closes the message store
func (mss *MessageStoreSQL) Close() error {
	if mss.db != nil {
		return mss.db.Close()
	}
	return nil
}

// Get retrieves a message by its CID
func (mss *MessageStoreSQL) Get(tenant, cidString string, options *MessageStoreOptions) (GenericMessage, error) {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return nil, options.Signal.Err()
		default:
		}
	}

	var encodedMessageBytes []byte
	var encodedData sql.NullString
	query := `SELECT encoded_message_bytes, encoded_data FROM message_store WHERE tenant = $1 AND message_cid = $2`
	err := mss.db.QueryRow(query, tenant, cidString).Scan(&encodedMessageBytes, &encodedData)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	var message GenericMessage
	err = cbor.Unmarshal(encodedMessageBytes, &message)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// If we have encoded data, add it back to the message
	if encodedData.Valid {
		messageMap, ok := message.(map[string]interface{})
		if ok {
			messageMap["encodedData"] = encodedData.String
		}
	}

	return message, nil
}

// Put stores a new message in the store
func (mss *MessageStoreSQL) Put(tenant string, message GenericMessage, indexes KeyValues, options *MessageStoreOptions) error {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return options.Signal.Err()
		default:
		}
	}

	// Extract and remove encodedData if present
	var encodedData sql.NullString
	if messageMap, ok := message.(map[string]interface{}); ok {
		if data, exists := messageMap["encodedData"]; exists {
			if dataStr, ok := data.(string); ok {
				encodedData = sql.NullString{String: dataStr, Valid: true}
			}
			delete(messageMap, "encodedData")
		}
	}

	// Encode the message
	encodedMessage, err := cbornode.WrapObject(message, multihash.SHA2_256, -1)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	messageCid := encodedMessage.Cid()
	encodedMessageBytes := encodedMessage.RawData()

	// Begin transaction
	tx, err := mss.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Prepare the insert statement with all index columns
	query := `
		INSERT INTO message_store (
			tenant, message_cid, encoded_message_bytes, encoded_data,
			interface, method, schema, data_cid, data_size,
			date_created, message_timestamp, data_format,
			is_latest_base_state, published, author, record_id,
			entry_id, date_published, latest, protocol,
			date_expires, description, granted_to, granted_by,
			granted_for, permissions_request_id, attester,
			protocol_path, recipient, context_id, parent_id,
			permissions_grant_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15, $16, $17, $18, $19,
			$20, $21, $22, $23, $24, $25, $26, $27, $28,
			$29, $30, $31, $32
		)`

	// Create args slice with all values
	args := []interface{}{
		tenant, messageCid.String(), encodedMessageBytes, encodedData,
		indexes["interface"], indexes["method"], indexes["schema"],
		indexes["dataCid"], indexes["dataSize"], indexes["dateCreated"],
		indexes["messageTimestamp"], indexes["dataFormat"],
		indexes["isLatestBaseState"], indexes["published"],
		indexes["author"], indexes["recordId"], indexes["entryId"],
		indexes["datePublished"], indexes["latest"], indexes["protocol"],
		indexes["dateExpires"], indexes["description"],
		indexes["grantedTo"], indexes["grantedBy"], indexes["grantedFor"],
		indexes["permissionsRequestId"], indexes["attester"],
		indexes["protocolPath"], indexes["recipient"],
		indexes["contextId"], indexes["parentId"],
		indexes["permissionsGrantId"],
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert message: %w", err)
	}

	return tx.Commit()
}

// Delete removes a message from the store
func (mss *MessageStoreSQL) Delete(tenant, cidString string, options *MessageStoreOptions) error {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return options.Signal.Err()
		default:
		}
	}

	query := `DELETE FROM message_store WHERE tenant = $1 AND message_cid = $2`
	result, err := mss.db.Exec(query, tenant, cidString)
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

// Query retrieves messages based on filters and sorting options
func (mss *MessageStoreSQL) Query(tenant string, filters []DataFilter, messageSort *MessageSort, pagination *Pagination, options *MessageStoreOptions) ([]GenericMessage, string, error) {
	if options != nil && options.Signal != nil {
		select {
		case <-options.Signal.Done():
			return nil, "", options.Signal.Err()
		default:
		}
	}

	// Build the base query
	query := `SELECT encoded_message_bytes, encoded_data FROM message_store WHERE tenant = $1`
	args := []interface{}{tenant}
	paramCount := 1

	// Add filters
	for _, filter := range filters {
		paramCount++
		query += fmt.Sprintf(" AND %s %s $%d", filter.Property, filter.Operator, paramCount)
		args = append(args, filter.Value)
	}

	// Add sorting
	if messageSort != nil && messageSort.Property != "" {
		direction := messageSort.Direction
		if direction == "" {
			direction = "ASC"
		}
		query += fmt.Sprintf(" ORDER BY %s %s", messageSort.Property, direction)
	}

	// Add pagination
	if pagination != nil && pagination.Limit > 0 {
		paramCount++
		query += fmt.Sprintf(" LIMIT $%d", paramCount)
		args = append(args, pagination.Limit)

		if pagination.Cursor != "" {
			paramCount++
			query += fmt.Sprintf(" OFFSET $%d", paramCount)
			offset, err := decodeCursor(pagination.Cursor)
			if err != nil {
				return nil, "", fmt.Errorf("invalid cursor: %w", err)
			}
			args = append(args, offset)
		}
	}

	// Execute query
	rows, err := mss.db.Query(query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var messages []GenericMessage
	for rows.Next() {
		var encodedMessageBytes []byte
		var encodedData sql.NullString

		err := rows.Scan(&encodedMessageBytes, &encodedData)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan row: %w", err)
		}

		var message GenericMessage
		err = cbor.Unmarshal(encodedMessageBytes, &message)
		if err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal message: %w", err)
		}

		// If we have encoded data, add it back to the message
		if encodedData.Valid {
			if messageMap, ok := message.(map[string]interface{}); ok {
				messageMap["encodedData"] = encodedData.String
			}
		}

		messages = append(messages, message)
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating rows: %w", err)
	}

	// Generate next cursor
	var nextCursor string
	if pagination != nil && len(messages) == pagination.Limit {
		offset := 0
		if pagination.Cursor != "" {
			offset, _ = decodeCursor(pagination.Cursor)
		}
		nextCursor = encodeCursor(offset + len(messages))
	}

	return messages, nextCursor, nil
}

// Helper functions for cursor encoding/decoding
func encodeCursor(offset int) string {
	return fmt.Sprintf("%d", offset)
}

func decodeCursor(cursor string) (int, error) {
	var offset int
	_, err := fmt.Sscanf(cursor, "%d", &offset)
	if err != nil {
		return 0, err
	}
	return offset, nil
}

// Clear removes all messages from the store
func (mss *MessageStoreSQL) Clear() error {
	_, err := mss.db.Exec("DELETE FROM message_store")
	if err != nil {
		return fmt.Errorf("failed to clear message store: %w", err)
	}
	return nil
}
