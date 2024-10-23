package store

import (
	"testing"
)

func TestCreateBlockStore(t *testing.T) {
	bs, err := NewBlockstoreLevel("data/test-blockstore")
	if err != nil {
		t.Fatalf("Failed to create BlockStore: %v", err)
	}
	if bs == nil {
		t.Fatal("BlockStore is nil")
	}
	// Add more specific tests based on BlockStore implementation
}

func TestCreateDataStore(t *testing.T) {
	ds, err := NewDataStoreLevel(DataStoreLevelConfig{
		BlockstoreLocation: "data/test-datastore",
	})
	if err != nil {
		t.Fatalf("Failed to create DataStore: %v", err)
	}
	if ds == nil {
		t.Fatal("DataStore is nil")
	}
	// Add more specific tests based on DataStore implementation
}

func TestCreateMessageStore(t *testing.T) {
	ms, err := NewMessageStoreLevel(MessageStoreLevelConfig{
		BlockstoreLocation: "data/test-messagestore",
		IndexLocation:      "data/test-index",
	})
	if err != nil {
		t.Fatalf("Failed to create MessageStore: %v", err)
	}
	if ms == nil {
		t.Fatal("MessageStore is nil")
	}
	// Add more specific tests based on MessageStore implementation
}
