package store

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEventLogSQLWithoutCleanup(t *testing.T) *EventLogSQL {
	config := getTestConfig()
	store, err := NewEventLogSQL(config)
	require.NoError(t, err)

	db, err := models.GetDB(config.DBConfig)
	require.NoError(t, err)
	store.db = db

	return store
}

func cleanupTestEventLogSQL(t *testing.T, store *EventLogSQL) {
	require.NoError(t, store.db.Exec("DELETE FROM event_logs").Error)
}

func setupTestEventLogSQL(t *testing.T) *EventLogSQL {
	store := setupTestEventLogSQLWithoutCleanup(t)
	cleanupTestEventLogSQL(t, store)
	return store
}

func TestEventLogSQL_OpenClose(t *testing.T) {
	store := setupTestEventLogSQL(t)

	err := store.Open()
	assert.NoError(t, err)

	err = store.Close()
	assert.NoError(t, err)

	assert.Nil(t, store.db)
}

func TestEventLogSQL_Append(t *testing.T) {
	store := setupTestEventLogSQLWithoutCleanup(t)
	defer cleanupTestEventLogSQL(t, store)

	now := time.Now().UTC().Format(time.RFC3339)
	metadata := map[string]interface{}{
		"test": "metadata",
		"foo":  "bar",
	}
	metadataBytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	indexes := KeyValues{
		"interface":            "test-interface",
		"method":               "test-method",
		"eventType":            "test-event",
		"schema":               "test-schema",
		"dataCid":              "test-dataCid",
		"dataSize":             "1024",
		"dateCreated":          now,
		"messageTimestamp":     now,
		"dataFormat":           "application/json",
		"isLatestBaseState":    "true",
		"published":            "true",
		"author":               "test-author",
		"recordId":             "test-recordId",
		"entryId":              "test-entryId",
		"datePublished":        now,
		"latest":               "true",
		"protocol":             "test-protocol",
		"dateExpires":          now,
		"description":          "test description",
		"grantedTo":            "test-grantedTo",
		"grantedBy":            "test-grantedBy",
		"grantedFor":           "test-grantedFor",
		"permissionsRequestId": "test-permReqId",
		"attester":             "test-attester",
		"protocolPath":         "test-protocolPath",
		"recipient":            "test-recipient",
		"contextId":            "test-contextId",
		"parentId":             "test-parentId",
		"permissionsGrantId":   "test-permGrantId",
		"metadata":             string(metadataBytes),
		"timestamp":            now,
	}

	err = store.Append("test-tenant", "test-cid", indexes)
	assert.NoError(t, err)

	// Verify event was stored with all fields
	var event models.EventLog
	err = store.db.Where("tenant = ? AND message_cid = ?", "test-tenant", "test-cid").First(&event).Error
	require.NoError(t, err)

	// Verify core fields
	assert.Equal(t, "test-interface", event.Interface)
	assert.Equal(t, "test-method", event.Method)
	assert.Equal(t, "test-event", event.EventType)
	assert.Equal(t, "test-schema", event.Schema)
	assert.Equal(t, "test-dataCid", event.DataCid)
	assert.Equal(t, "1024", event.DataSize)

	// Verify timestamps
	assert.Equal(t, now, event.DateCreated)
	assert.Equal(t, now, event.MessageTimestamp)
	assert.Equal(t, now, event.DatePublished)
	assert.Equal(t, now, event.DateExpires)

	// Verify metadata
	var expectedMetadata map[string]interface{}
	var actualMetadata map[string]interface{}
	err = json.Unmarshal(metadataBytes, &expectedMetadata)
	require.NoError(t, err)
	err = json.Unmarshal(event.Metadata, &actualMetadata)
	require.NoError(t, err)
	assert.Equal(t, expectedMetadata, actualMetadata)

	// Verify additional fields
	assert.Equal(t, "test-author", event.Author)
	assert.Equal(t, "test-recordId", event.RecordId)
	assert.Equal(t, "test-protocol", event.Protocol)
	assert.Equal(t, "test-recipient", event.Recipient)
	assert.Equal(t, "test-parentId", event.ParentId)
}

func TestEventLogSQL_QueryEvents(t *testing.T) {
	store := setupTestEventLogSQLWithoutCleanup(t)
	defer cleanupTestEventLogSQL(t, store)

	// Add test events
	events := []struct {
		tenant     string
		messageCid string
		eventType  string
	}{
		{"tenant1", "cid1", "type1"},
		{"tenant1", "cid2", "type2"},
		{"tenant2", "cid3", "type1"},
	}

	metadata := map[string]interface{}{
		"test": "metadata",
	}
	metadataBytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	for _, e := range events {
		err := store.Append(e.tenant, e.messageCid, KeyValues{
			"eventType": e.eventType,
			"metadata":  string(metadataBytes),
		})
		require.NoError(t, err)
	}

	// Test query with filter
	filters := []DataFilter{
		{Property: "event_type", Operator: "=", Value: "type1"},
	}
	results, err := store.QueryEvents("tenant1", filters, nil)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "cid1", results[0])

	// Test cursor-based pagination
	results, err = store.QueryEvents("tenant1", nil, &EventOptions{Cursor: "cid1"})
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "cid2", results[0])
}

func TestEventLogSQL_DeleteEventsByCid(t *testing.T) {
	store := setupTestEventLogSQLWithoutCleanup(t)
	defer cleanupTestEventLogSQL(t, store)

	// Initialize metadata
	metadata := map[string]interface{}{
		"test": "metadata",
	}
	metadataBytes, err := json.Marshal(metadata)
	require.NoError(t, err)

	// Add test event
	err = store.Append("test-tenant", "test-cid", KeyValues{
		"eventType": "test-event",
		"metadata":  string(metadataBytes),
	})
	require.NoError(t, err)

	// Delete event
	err = store.DeleteEventsByCid("test-tenant", []string{"test-cid"})
	assert.NoError(t, err)

	// Verify deletion
	var count int64
	err = store.db.Model(&models.EventLog{}).Where("tenant = ? AND message_cid = ?", "test-tenant", "test-cid").Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count)
}
