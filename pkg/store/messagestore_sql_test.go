package store

import (
	"testing"
	"time"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/config"
	"github.com/abaxxtech/abaxx-id-go/pkg/store/models"
	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestConfig() MessageStoreSQLConfig {
	return MessageStoreSQLConfig{
		DBConfig: config.NewDefaultConfig(),
	}
}

func setupTestGormStore(t *testing.T) *GormMessageStore {
	config := getTestConfig()
	store, err := NewMessageStoreSQL(config)
	require.NoError(t, err)

	// Get the singleton DB instance
	db, err := models.GetDB(config.DBConfig)
	require.NoError(t, err)
	store.db = db

	// Clear existing data
	err = db.Exec("DELETE FROM message_store").Error
	require.NoError(t, err)

	return store
}

func setupTestGormStoreWithoutCleanup(t *testing.T) *GormMessageStore {
	config := getTestConfig()
	store, err := NewMessageStoreSQL(config)
	require.NoError(t, err)

	// Get the singleton DB instance
	db, err := models.GetDB(config.DBConfig)
	require.NoError(t, err)
	store.db = db

	return store
}

func cleanupTestGormStore(t *testing.T, store *GormMessageStore) {
	if store != nil && store.db != nil {
		err := store.db.Exec("DELETE FROM message_store").Error
		require.NoError(t, err)
	}
}

func TestDatabaseSetup(t *testing.T) {
	store := setupTestGormStoreWithoutCleanup(t)

	// First verify the database is empty
	var initialCount int64
	err := store.db.Model(&models.MessageStore{}).Count(&initialCount).Error
	require.NoError(t, err)
	require.Equal(t, int64(0), initialCount, "Database should be empty initially")

	// Create test message with all required fields from the model
	messageBytes := []byte(`{"cid":"test-cid","content":"test message"}`)
	message := map[string]interface{}{
		"cid":                 "test-cid",
		"content":             "test message",
		"encodedMessageBytes": messageBytes,
	}

	indexes := KeyValues{
		"interface":         "test-interface",
		"method":            "test-method",
		"dateCreated":       time.Now().UTC().Format(time.RFC3339),
		"schema":            "test-schema",
		"dataCid":           "test-dataCid",
		"dataSize":          "1024",
		"protocol":          "test-protocol",
		"recordId":          "test-recordId",
		"entryId":           "test-entryId",
		"datePublished":     time.Now().UTC().Format(time.RFC3339),
		"encodedData":       string(messageBytes),
		"messageTimestamp":  time.Now().UTC().Format(time.RFC3339),
		"dataFormat":        "application/json",
		"isLatestBaseState": "true",
		"published":         "true",
		"latest":            "true",
	}

	// Test Put
	err = store.Put("test-tenant", message, indexes, nil)
	require.NoError(t, err)

	// Verify the message was stored
	var count int64
	err = store.db.Model(&models.MessageStore{}).Count(&count).Error
	require.NoError(t, err)
	require.Equal(t, int64(1), count, "One record should be in database")

	// Query the message directly from the database - without specifying the CID
	var storedMessage models.MessageStore
	err = store.db.Where("tenant = ?", "test-tenant").First(&storedMessage).Error
	require.NoError(t, err)

	// Print the full stored message for debugging
	t.Logf("Stored message: %+v", storedMessage)

	// Verify key fields are present
	require.Equal(t, "test-tenant", storedMessage.Tenant)
	require.Equal(t, "test-interface", storedMessage.Interface)
	require.Equal(t, "test-method", storedMessage.Method)
	require.Equal(t, "test-schema", storedMessage.Schema)
	require.Equal(t, "test-dataCid", storedMessage.DataCid)
	require.Equal(t, "1024", storedMessage.DataSize)
	require.Equal(t, "test-protocol", storedMessage.Protocol)
	require.Equal(t, "true", storedMessage.Latest)
	require.Equal(t, "true", storedMessage.Published)
	require.Equal(t, "true", storedMessage.IsLatestBaseState)

	// Store the generated CID for future reference
	generatedCID := storedMessage.MessageCid
	t.Logf("Generated CID: %s", generatedCID)

	// Use the version without cleanup for the persistence check
	newStore := setupTestGormStoreWithoutCleanup(t)
	retrieved, err := newStore.Get("test-tenant", generatedCID)
	require.NoError(t, err)
	require.NotNil(t, retrieved)
}

func TestGormMessageStore_OpenClose(t *testing.T) {
	store := setupTestGormStoreWithoutCleanup(t)
	// defer cleanupTestGormStore(t, store)

	// Since we're using singleton DB, Open() should return nil error
	// as the connection is already established
	err := store.Open()
	assert.NoError(t, err)

	// Close() should not actually close the singleton connection
	err = store.Close()
	assert.NoError(t, err)

	// Verify DB is still usable through the singleton
	db, err := models.GetDB(store.config.DBConfig)
	assert.NoError(t, err)
	err = db.Exec("SELECT 1").Error
	assert.NoError(t, err)
}

func TestGormMessageStore_PutGet(t *testing.T) {
	store := setupTestGormStoreWithoutCleanup(t)
	defer cleanupTestGormStore(t, store)

	// Create test message with required fields
	messageBytes := []byte(`{"cid":"test-cid","content":"test message"}`)
	message := map[string]interface{}{
		"cid":                 "test-cid",
		"content":             "test message",
		"encodedMessageBytes": messageBytes,
	}

	indexes := KeyValues{
		"interface":     "test-interface",
		"method":        "test-method",
		"dateCreated":   time.Now().String(),
		"encodedData":   "test-data",
		"schema":        "test-schema",
		"dataCid":       "test-cid",
		"dataSize":      "1024",
		"protocol":      "test-protocol",
		"recordId":      "test-recordId",
		"entryId":       "test-entryId",
		"datePublished": time.Now().String(),
	}

	// Test Put
	err := store.Put("test-tenant", message, indexes, nil)
	assert.NoError(t, err)

	// Query the message directly from the database to get the generated CID
	var storedMessage models.MessageStore
	err = store.db.Where("tenant = ?", "test-tenant").First(&storedMessage).Error
	require.NoError(t, err)

	// Store the generated CID
	generatedCID := storedMessage.MessageCid

	// Decode the generated CID to cid.Cid type
	decodedCID, err := cid.Decode(generatedCID)
	require.NoError(t, err)

	// Test Get - use the decoded CID
	retrieved, err := store.Get("test-tenant", decodedCID.String())
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	// Verify the retrieved message
	retrievedMap, ok := retrieved.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, message["cid"], retrievedMap["cid"])
	assert.Equal(t, message["content"], retrievedMap["content"])

	// Test Get non-existent
	retrieved, err = store.Get("test-tenant", "non-existent")
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestGormMessageStore_Query(t *testing.T) {
	store := setupTestGormStoreWithoutCleanup(t)
	defer cleanupTestGormStore(t, store)

	// Insert multiple test messages
	messages := []struct {
		cid        string
		content    string
		interface_ string
		method     string
	}{
		{"cid-1", "content 1", "interface-1", "method-1"},
		{"cid-2", "content 2", "interface-1", "method-2"},
		{"cid-3", "content 3", "interface-2", "method-1"},
		{"cid-4", "content 4", "interface-1", "method-1"},
	}

	for _, msg := range messages {
		message := map[string]interface{}{
			"cid":     msg.cid,
			"content": msg.content,
		}

		indexes := KeyValues{
			"interface":   msg.interface_,
			"method":      msg.method,
			"dateCreated": time.Now().Format(time.RFC3339),
		}

		err := store.Put("test-tenant", message, indexes, nil)
		require.NoError(t, err)
	}

	// Verify test data is inserted
	var count int64
	err := store.db.Model(&models.MessageStore{}).Count(&count).Error
	require.NoError(t, err)
	//assert.Equal(t, int64(len(messages)), count)

	// Test query with filters
	filters := []DataFilter{
		{Property: "interface", Operator: "=", Value: "interface-1"},
	}

	sort := &MessageSort{
		Property:  "date_created",
		Direction: "DESC",
	}

	pagination := &Pagination{
		Limit: 10,
	}

	result, cursor, err := store.Query("test-tenant", filters, sort, pagination, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	//assert.NotEmpty(t, cursor)

	// Verify the number of returned messages
	assert.Len(t, result, 3)

	// Verify the returned messages have the expected interface
	// for _, msg := range result {
	// 	msgStore, ok := msg.(*models.MessageStore)
	// 	require.True(t, ok, "Expected message to be of type *models.MessageStore")
	// 	assert.Equal(t, "interface-1", msgStore.Interface)
	// }

	// Test query with multiple filters
	filters = append(filters, DataFilter{Property: "method", Operator: "=", Value: "method-1"})
	result, cursor, err = store.Query("test-tenant", filters, sort, pagination, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	//assert.NotEmpty(t, cursor)

	// Verify the number of returned messages
	assert.Len(t, result, 2)

	// Verify the returned messages have the expected interface and method
	// for _, msg := range result {
	// 	msgStore, ok := msg.(*models.MessageStore)
	// 	require.True(t, ok, "Expected message to be of type *models.MessageStore")
	// 	assert.Equal(t, "interface-1", msgStore.Interface)
	// 	assert.Equal(t, "method-1", msgStore.Method)
	// }

	// Test pagination
	pagination.Limit = 1
	result, cursor, err = store.Query("test-tenant", nil, sort, pagination, nil)
	require.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.NotEmpty(t, cursor)
	assert.Len(t, result, 1)
}

func TestGormMessageStore_Delete(t *testing.T) {
	store := setupTestGormStoreWithoutCleanup(t)
	defer cleanupTestGormStore(t, store)

	// Create test message
	messageBytes := []byte(`{"cid":"test-cid","content":"test message"}`)
	message := map[string]interface{}{
		"content":             "test message",
		"encodedMessageBytes": messageBytes,
	}

	indexes := KeyValues{
		"interface":   "test-interface",
		"method":      "test-method",
		"dateCreated": time.Now().String(),
	}

	// Put a message
	err := store.Put("test-tenant", message, indexes, nil)
	assert.NoError(t, err)

	// Query the message directly from the database to get the generated CID
	var storedMessage models.MessageStore
	err = store.db.Where("tenant = ?", "test-tenant").First(&storedMessage).Error
	require.NoError(t, err)

	// Store the generated CID
	generatedCID := storedMessage.MessageCid

	// Delete the message using the generated CID
	err = store.Delete("test-tenant", generatedCID, nil)
	assert.NoError(t, err)

	// Verify it's soft deleted by checking deleted_at is set
	var deletedMessage models.MessageStore
	err = store.db.Unscoped().Where("message_cid = ?", generatedCID).First(&deletedMessage).Error
	require.NoError(t, err)
	assert.NotNil(t, deletedMessage.DeletedAt, "Message should be soft deleted")

	// Verify Get returns nil for deleted message
	retrieved, err := store.Get("test-tenant", generatedCID)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)

	// Try to delete non-existent message
	err = store.Delete("test-tenant", "non-existent", nil)
	assert.Error(t, err)
}

func TestGormMessageStore_MultiTenant(t *testing.T) {
	store := setupTestGormStore(t)
	defer cleanupTestGormStore(t, store)

	messageBytes := []byte(`{"cid":"test-cid","content":"test message"}`)
	message := map[string]interface{}{
		"content":             "test message",
		"encodedMessageBytes": messageBytes,
	}

	indexes := KeyValues{
		"interface":   "test-interface",
		"method":      "test-method",
		"dateCreated": time.Now().String(),
	}

	// Put same message for different tenants
	err := store.Put("tenant1", message, indexes, nil)
	assert.NoError(t, err)
	err = store.Put("tenant2", message, indexes, nil)
	assert.NoError(t, err)

	// Get the generated CIDs for both tenants
	var tenant1Message models.MessageStore
	err = store.db.Where("tenant = ?", "tenant1").First(&tenant1Message).Error
	require.NoError(t, err)
	tenant1CID := tenant1Message.MessageCid

	var tenant2Message models.MessageStore
	err = store.db.Where("tenant = ?", "tenant2").First(&tenant2Message).Error
	require.NoError(t, err)
	tenant2CID := tenant2Message.MessageCid

	// Verify messages exist for both tenants using generated CIDs
	retrieved1, err := store.Get("tenant1", tenant1CID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved1)

	retrieved2, err := store.Get("tenant2", tenant2CID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved2)

	// Delete message for tenant1
	err = store.Delete("tenant1", tenant1CID, nil)
	assert.NoError(t, err)

	// Verify message still exists for tenant2
	retrieved2, err = store.Get("tenant2", tenant2CID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved2)
}
