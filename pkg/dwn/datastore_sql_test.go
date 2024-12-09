package dwn

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestDBConfig() SQLStoreConfig {
	// Use environment variables or default test database
	dbHost := os.Getenv("TEST_DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("TEST_DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbName := os.Getenv("TEST_DB_NAME")
	if dbName == "" {
		dbName = "dwn"
	}
	dbUser := os.Getenv("TEST_DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("TEST_DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}

	return SQLStoreConfig{
		DriverName: "postgres",
		DataSourceName: fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			dbUser, dbPassword, dbHost, dbPort, dbName,
		),
	}
}

func setupTestStore(t *testing.T) *SQLStore {
	config := getTestDBConfig()
	store, err := NewSQLStore(config)
	require.NoError(t, err)
	require.NoError(t, store.Open())
	require.NoError(t, store.Clear()) // Start with a clean state
	return store
}

func cleanupTestStore(t *testing.T, store *SQLStore) {
	require.NoError(t, store.Clear())
	require.NoError(t, store.Close())
}

func TestSQLStore_OpenClose(t *testing.T) {
	store := setupTestStore(t)
	defer cleanupTestStore(t, store)

	// Test double open
	err := store.Open()
	assert.Error(t, err)

	// Test close
	err = store.Close()
	assert.NoError(t, err)

	// Test double close
	err = store.Close()
	assert.NoError(t, err)
}

func TestSQLStore_PutGet(t *testing.T) {
	store := setupTestStore(t)
	defer cleanupTestStore(t, store)

	tenant := Tenant("test-tenant")
	message := map[string]interface{}{
		"cid":     "test-cid",
		"content": "test message",
	}
	indexes := IndexableKeyValues{
		"interface": S("test-interface"),
		"method":    S("test-method"),
	}

	// Test Put
	err := store.Put(tenant, message, indexes)
	assert.NoError(t, err)

	// Test Get
	retrieved, err := store.Get(tenant, MessageCid("test-cid"))
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	retrievedMap, ok := retrieved.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, message["cid"], retrievedMap["cid"])
	assert.Equal(t, message["content"], retrievedMap["content"])

	// Test Get non-existent message
	retrieved, err = store.Get(tenant, MessageCid("non-existent"))
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestSQLStore_Delete(t *testing.T) {
	store := setupTestStore(t)
	defer cleanupTestStore(t, store)

	tenant := Tenant("test-tenant")
	message := map[string]interface{}{
		"cid":     "test-cid",
		"content": "test message",
	}
	indexes := IndexableKeyValues{}

	// Put a message
	err := store.Put(tenant, message, indexes)
	assert.NoError(t, err)

	// Delete the message
	err = store.Delete(tenant, MessageCid("test-cid"))
	assert.NoError(t, err)

	// Verify it's deleted
	retrieved, err := store.Get(tenant, MessageCid("test-cid"))
	assert.NoError(t, err)
	assert.Nil(t, retrieved)

	// Try to delete non-existent message
	err = store.Delete(tenant, MessageCid("non-existent"))
	assert.Error(t, err)
}

func TestSQLStore_Query(t *testing.T) {
	store := setupTestStore(t)
	defer cleanupTestStore(t, store)

	tenant := Tenant("test-tenant")
	messages := []map[string]interface{}{
		{
			"cid":     "msg1",
			"content": "message 1",
			"type":    "type1",
		},
		{
			"cid":     "msg2",
			"content": "message 2",
			"type":    "type2",
		},
	}

	// Put test messages
	for _, msg := range messages {
		err := store.Put(tenant, msg, IndexableKeyValues{
			"type": S(msg["type"].(string)),
		})
		assert.NoError(t, err)
	}

	// Test query with filter
	filters := []Filter{
		&testFilter{property: "type", value: EqualFilter{EqualTo: S("type1")}},
	}
	err := store.Query(tenant, filters, MessageSort{}, Pagination{})
	assert.NoError(t, err)

	// Test query with pagination
	pagination := Pagination{
		Limit:  1,
		Offset: 0,
	}
	err = store.Query(tenant, nil, MessageSort{}, pagination)
	assert.NoError(t, err)
}

// Helper test filter implementation
type testFilter struct {
	property string
	value    FilterValue
}

func (f *testFilter) Property() string {
	return f.property
}

func (f *testFilter) Value() FilterValue {
	return f.value
}

func TestSQLStore_MultiTenant(t *testing.T) {
	store := setupTestStore(t)
	defer cleanupTestStore(t, store)

	tenant1 := Tenant("tenant1")
	tenant2 := Tenant("tenant2")
	message := map[string]interface{}{
		"cid":     "test-cid",
		"content": "test message",
	}
	indexes := IndexableKeyValues{}

	// Put same message for different tenants
	err := store.Put(tenant1, message, indexes)
	assert.NoError(t, err)
	err = store.Put(tenant2, message, indexes)
	assert.NoError(t, err)

	// Get message for tenant1
	retrieved1, err := store.Get(tenant1, MessageCid("test-cid"))
	assert.NoError(t, err)
	assert.NotNil(t, retrieved1)

	// Get message for tenant2
	retrieved2, err := store.Get(tenant2, MessageCid("test-cid"))
	assert.NoError(t, err)
	assert.NotNil(t, retrieved2)

	// Delete message for tenant1
	err = store.Delete(tenant1, MessageCid("test-cid"))
	assert.NoError(t, err)

	// Verify message still exists for tenant2
	retrieved2, err = store.Get(tenant2, MessageCid("test-cid"))
	assert.NoError(t, err)
	assert.NotNil(t, retrieved2)
}

func TestSQLStore_UpdateMessage(t *testing.T) {
	store := setupTestStore(t)
	defer cleanupTestStore(t, store)

	tenant := Tenant("test-tenant")
	message1 := map[string]interface{}{
		"cid":     "test-cid",
		"content": "original content",
	}
	message2 := map[string]interface{}{
		"cid":     "test-cid",
		"content": "updated content",
	}
	indexes := IndexableKeyValues{}

	// Put original message
	err := store.Put(tenant, message1, indexes)
	assert.NoError(t, err)

	// Update message
	err = store.Put(tenant, message2, indexes)
	assert.NoError(t, err)

	// Get message and verify it's updated
	retrieved, err := store.Get(tenant, MessageCid("test-cid"))
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)

	retrievedMap, ok := retrieved.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, message2["content"], retrievedMap["content"])
}
