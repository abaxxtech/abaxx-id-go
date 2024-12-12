package models

import (
	"testing"
	"time"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestData represents sample data for testing
type TestData struct {
	MessageStore       []MessageStore
	DataStore          []DataStore
	DataStoreReference []DataStoreReference
	EventLog           []EventLog
}

func setupTestDB(t *testing.T) (*gorm.DB, TestData) {
	// Try to get DB connection with default config
	db, err := GetDB(config.NewDefaultConfig())
	if err != nil {
		t.Logf("Warning: Could not connect to database: %v", err)
		return nil, TestData{}
	}

	// Create test data only if we have a valid connection
	testData := TestData{
		MessageStore: []MessageStore{
			{
				Tenant:               "test-tenant-1",
				MessageCid:           "test-cid-1",
				EncodedMessageBytes:  []byte("test-message-1"),
				Interface:            "test-interface",
				Method:               "test-method",
				DateCreated:          time.Now().String(),
				Schema:               "test-schema",
				DataCid:              "test-dataCid-1",
				DataSize:             "1024",
				Protocol:             "test-protocol",
				RecordId:             "test-recordId",
				EntryId:              "test-entryId",
				DatePublished:        time.Now().String(),
				Latest:               "true",
				DataFormat:           "application/json",
				IsLatestBaseState:    "true",
				Published:            "true",
				Author:               "test-author",
				GrantedTo:            "test-grantedTo",
				GrantedBy:            "test-grantedBy",
				GrantedFor:           "test-grantedFor",
				PermissionsRequestId: "test-permissionsRequestId",
				Attester:             "test-attester",
				ProtocolPath:         "test-protocolPath",
				Recipient:            "test-recipient",
				ContextId:            "test-contextId",
				ParentId:             "test-parentId",
				PermissionsGrantId:   "test-permissionsGrantId",
				Description:          "test-description",
				DateExpires:          time.Now().String(),
			},
			{
				Tenant:              "test-tenant-2",
				MessageCid:          "test-cid-2",
				EncodedMessageBytes: []byte("test-message-2"),
				Interface:           "test-interface",
				Method:              "test-method",
				DateCreated:         time.Now().String(),
			},
		},
		DataStore: []DataStore{
			{
				Tenant:      "test-tenant-1",
				DataCid:     "test-data-cid-1",
				EncodedData: []byte("test-data-1"),
				DataFormat:  "application/json",
				Schema:      "test-schema",
				DataSize:    1024,
				DateCreated: time.Now().String(),
			},
		},
		DataStoreReference: []DataStoreReference{
			{
				Tenant:      "test-tenant-1",
				MessageCid:  "test-cid-1",
				DataCid:     "test-data-cid-1",
				DateCreated: time.Now().String(),
			},
		},
		EventLog: []EventLog{
			{
				Tenant:      "test-tenant-1",
				EventType:   "test-event",
				MessageCid:  "test-cid-1",
				DataCid:     "test-data-cid-1",
				Description: "test description",
				Metadata:    []byte(`{"test": "metadata"}`),
				Timestamp:   time.Now().String(),
			},
		},
	}

	// Only try to insert test data if we have a valid connection
	if db != nil {
		for _, msg := range testData.MessageStore {
			if err := db.Create(&msg).Error; err != nil {
				t.Logf("Warning: Failed to create test message: %v", err)
			}
		}
		for _, data := range testData.DataStore {
			require.NoError(t, db.Create(&data).Error)
		}
		for _, ref := range testData.DataStoreReference {
			require.NoError(t, db.Create(&ref).Error)
		}
		for _, event := range testData.EventLog {
			require.NoError(t, db.Create(&event).Error)
		}
	}

	return db, testData
}

func cleanupTestDB(t *testing.T, db *gorm.DB) {
	// Clean up in reverse order of dependencies
	require.NoError(t, db.Exec("DELETE FROM event_logs").Error)
	require.NoError(t, db.Exec("DELETE FROM data_store_references").Error)
	require.NoError(t, db.Exec("DELETE FROM data_stores").Error)
	require.NoError(t, db.Exec("DELETE FROM message_store").Error)
}

func TestGetDB_InvalidConfig(t *testing.T) {
	// Skip if no database connection is available
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	invalidConfig := config.DBConfig{
		Host:     "nonexistent",
		Port:     "5432",
		User:     "invalid",
		Password: "invalid",
		DBName:   "nonexistent",
		SSLMode:  "disable",
	}

	db, err := GetDB(invalidConfig)
	// We expect an error here, but we'll handle it gracefully
	assert.Error(t, err)
	assert.Nil(t, db)
	assert.Contains(t, err.Error(), "connect")
}

func TestGetDB_ValidConfig(t *testing.T) {
	// Skip if no database connection is available
	if testing.Short() {
		t.Skip("Skipping database integration test in short mode")
	}

	db, testData := setupTestDB(t)
	if db == nil {
		t.Skip("Database connection not available - skipping test")
		return
	}
	defer cleanupTestDB(t, db)

	// Test singleton behavior
	db2, err := GetDB(config.NewDefaultConfig())
	if err != nil {
		t.Skip("Database connection failed - skipping test")
		return
	}
	assert.Equal(t, db, db2, "Should return the same instance")

	// Verify test data
	var count int64
	db.Model(&MessageStore{}).Count(&count)
	assert.Equal(t, int64(len(testData.MessageStore)), count)
}

func TestGetDB_AutoMigration(t *testing.T) {
	db, _ := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Verify tables exist
	var tableNames []string
	err := db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Pluck("tablename", &tableNames).Error
	require.NoError(t, err)

	expectedTables := []string{
		"message_store",
		"data_stores",
		"data_store_references",
		"event_logs",
	}

	for _, tableName := range expectedTables {
		assert.Contains(t, tableNames, tableName)
	}
}
