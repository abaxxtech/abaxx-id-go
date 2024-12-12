package store

import (
	"testing"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/models"
	"github.com/stretchr/testify/require"
)

func setupTestDataStoreSQLWithoutCleanup(t *testing.T) *DataStoreSQL {
	config := getTestConfig()
	store, err := NewDataStoreSQL(config)
	require.NoError(t, err)

	db, err := models.GetDB(config.DBConfig)
	require.NoError(t, err)
	store.db = db

	return store
}

func cleanupTestDataStoreSQL(t *testing.T, store *DataStoreSQL) error {
	if store.db == nil {
		return nil
	}

	// Clear existing data
	err := store.db.Exec("DELETE FROM data_store_references").Error
	if err != nil {
		return err
	}

	err = store.db.Exec("DELETE FROM data_stores").Error
	if err != nil {
		return err
	}

	return nil
}

// func TestDataStoreSQL_OpenClose(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)

// 	err := store.Open()
// 	assert.NoError(t, err)

// 	err = store.Close()
// 	assert.NoError(t, err)

// 	assert.Nil(t, store.db)
// }

// func TestDataStoreSQL_PutGet(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)
// 	defer cleanupTestDataStoreSQL(t, store)

// 	testData := []byte("test data")
// 	dataReader := bytes.NewReader(testData)

// 	result, err := store.Put("test-tenant", "test-message", "test-cid", dataReader)
// 	require.NoError(t, err)
// 	assert.Equal(t, "test-cid", result.DataCid)
// 	assert.Equal(t, int64(len(testData)), result.DataSize)

// 	// Test Get
// 	resultCid, resultSize, resultReader, err := store.Get("test-tenant", "test-message", "test-cid")
// 	require.NoError(t, err)
// 	assert.Equal(t, "test-cid", resultCid)
// 	assert.Equal(t, int64(len(testData)), resultSize)

// 	resultData, err := io.ReadAll(resultReader)
// 	require.NoError(t, err)
// 	assert.Equal(t, testData, resultData)
// }

// func TestDataStoreSQL_Associate(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)
// 	defer func() {
// 		require.NoError(t, cleanupTestDataStoreSQL(t, store))
// 	}()

// 	// Put initial data
// 	testData := []byte("test data")
// 	dataReader := bytes.NewReader(testData)
// 	putResult, err := store.Put("test-tenant", "message-1", "test-cid", dataReader)
// 	require.NoError(t, err)
// 	require.NotNil(t, putResult)

// 	// Verify the data exists in data_stores before attempting to associate
// 	var dataStore models.DataStore
// 	result := store.db.Where(&models.DataStore{
// 		Tenant:  "test-tenant",
// 		DataCid: "test-cid",
// 	}).First(&dataStore)
// 	require.NoError(t, result.Error)

// 	// Associate with another message
// 	associateResult, err := store.Associate("test-tenant", "message-2", "test-cid")
// 	require.NoError(t, err)
// 	require.NotNil(t, associateResult)
// 	assert.Equal(t, "test-cid", associateResult.DataCid)
// 	assert.Equal(t, int64(len(testData)), associateResult.DataSize)

// 	// Verify both references exist
// 	var refs []models.DataStoreReference
// 	err = store.db.Where(&models.DataStoreReference{
// 		Tenant:  "test-tenant",
// 		DataCid: "test-cid",
// 	}).Find(&refs).Error
// 	require.NoError(t, err)
// 	assert.Equal(t, 2, len(refs))
// }

// func TestDataStoreSQL_Delete(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)
// 	defer func() {
// 		require.NoError(t, cleanupTestDataStoreSQL(t, store))
// 	}()

// 	// Put data and associate with two messages
// 	testData := []byte("test data")
// 	dataReader := bytes.NewReader(testData)
// 	putResult, err := store.Put("test-tenant", "message-1", "test-cid", dataReader)
// 	require.NoError(t, err)
// 	require.NotNil(t, putResult)

// 	// Associate with second message
// 	associateResult, err := store.Associate("test-tenant", "message-2", "test-cid")
// 	require.NoError(t, err)
// 	require.NotNil(t, associateResult)

// 	// Verify both references exist
// 	var refs []models.DataStoreReference
// 	err = store.db.Where(&models.DataStoreReference{
// 		Tenant:  "test-tenant",
// 		DataCid: "test-cid",
// 	}).Find(&refs).Error
// 	require.NoError(t, err)
// 	assert.Equal(t, 2, len(refs))

// 	// Delete first message's reference
// 	err = store.Delete("test-tenant", "message-1", "test-cid")
// 	require.NoError(t, err)

// 	// Verify second message can still access the data
// 	_, _, reader, err := store.Get("test-tenant", "message-2", "test-cid")
// 	require.NoError(t, err)
// 	require.NotNil(t, reader)
// 	data, err := io.ReadAll(reader)
// 	require.NoError(t, err)
// 	assert.Equal(t, testData, data)

// 	// Delete second message's reference
// 	err = store.Delete("test-tenant", "message-2", "test-cid")
// 	require.NoError(t, err)

// 	// Verify data is completely deleted
// 	_, _, reader, err = store.Get("test-tenant", "message-2", "test-cid")
// 	require.NoError(t, err)
// 	assert.Nil(t, reader)
// }

// func TestDataStoreSQL_Clear(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)

// 	// Put some test data
// 	testData := []byte("test data")
// 	dataReader := bytes.NewReader(testData)
// 	_, err := store.Put("test-tenant", "message-1", "test-cid", dataReader)
// 	require.NoError(t, err)

// 	// Clear all data
// 	err = store.Clear()
// 	require.NoError(t, err)

// 	// Verify data is cleared
// 	_, _, reader, err := store.Get("test-tenant", "message-1", "test-cid")
// 	require.NoError(t, err)
// 	assert.Nil(t, reader)
// }

// func TestDataStoreSQL_ErrorCases(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)

// 	// Test operations without opening DB
// 	store.db = nil

// 	// Test Put
// 	_, err := store.Put("test-tenant", "test-message", "test-cid", bytes.NewReader([]byte("test")))
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "connection to database not open")

// 	// Test Get
// 	_, _, reader, err := store.Get("test-tenant", "test-message", "test-cid")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "connection to database not open")
// 	assert.Nil(t, reader)

// 	// Test Associate
// 	_, err = store.Associate("test-tenant", "test-message", "test-cid")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "connection to database not open")

// 	// Test Delete
// 	err = store.Delete("test-tenant", "test-message", "test-cid")
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "connection to database not open")

// 	// Test Clear
// 	err = store.Clear()
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "connection to database not open")
// }

// func TestDataStoreSQL_NonExistentData(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)

// 	// Test Get non-existent data
// 	_, _, reader, err := store.Get("test-tenant", "non-existent", "test-cid")
// 	require.NoError(t, err)
// 	assert.Nil(t, reader)

// 	// Test Associate with non-existent data
// 	_, err = store.Associate("test-tenant", "test-message", "non-existent")
// 	require.NoError(t, err)
// 	assert.Nil(t, reader)

// 	// Test Delete non-existent data
// 	err = store.Delete("test-tenant", "non-existent", "test-cid")
// 	require.NoError(t, err)
// }

// func TestDataStoreSQL_MultiTenant(t *testing.T) {
// 	store := setupTestDataStoreSQLWithoutCleanup(t)

// 	// Put same data for different tenants
// 	testData := []byte("test data")
// 	dataReader1 := bytes.NewReader(testData)
// 	dataReader2 := bytes.NewReader(testData)

// 	_, err := store.Put("tenant-1", "message-1", "test-cid", dataReader1)
// 	require.NoError(t, err)
// 	_, err = store.Put("tenant-2", "message-1", "test-cid", dataReader2)
// 	require.NoError(t, err)

// 	// Verify data isolation between tenants
// 	_, _, reader1, err := store.Get("tenant-1", "message-1", "test-cid")
// 	require.NoError(t, err)
// 	require.NotNil(t, reader1)
// 	data1, err := io.ReadAll(reader1)
// 	require.NoError(t, err)
// 	assert.Equal(t, testData, data1)

// 	_, _, reader2, err := store.Get("tenant-2", "message-1", "test-cid")
// 	require.NoError(t, err)
// 	require.NotNil(t, reader2)
// 	data2, err := io.ReadAll(reader2)
// 	require.NoError(t, err)
// 	assert.Equal(t, testData, data2)

// 	// Delete data for tenant-1
// 	err = store.Delete("tenant-1", "message-1", "test-cid")
// 	require.NoError(t, err)

// 	// Verify tenant-2 data still exists
// 	_, _, reader, err := store.Get("tenant-2", "message-1", "test-cid")
// 	require.NoError(t, err)
// 	require.NotNil(t, reader)
// 	data, err := io.ReadAll(reader)
// 	require.NoError(t, err)
// 	assert.Equal(t, testData, data)
// }
