package dwn

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	cid "github.com/ipfs/go-cid"
	jwk "github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
)

// GetKeyById retrieves a key from the DID by its ID
const authSig = `
{ "signature": 
  {
    "payload": "eyJkZXNjcmlwdG9yQ2lkIjoiYmFmeXJlaWJpc3F2eG9sNHF6ZDVkem9wa3VlZHR3eTRrbnp1Z2dudWprbW15b25hNzV3dGxvd2k0ZGkifQ",
    "signatures": [
      {
	"protected": "eyJraWQiOiJkaWQ6aW9uOkVpQVRUNmZLQWc0cHJzQVE5VkE5ZGcxbmZHSTBLdUtfVTNzUUdrcFZvVDlNdkE6ZXlKa1pXeDBZU0k2ZXlKd1lYUmphR1Z6SWpwYmV5SmhZM1JwYjI0aU9pSnlaWEJzWVdObElpd2laRzlqZFcxbGJuUWlPbnNpY0hWaWJHbGpTMlY1Y3lJNlczc2lhV1FpT2lKa2QyNHRjMmxuSWl3aWNIVmliR2xqUzJWNVNuZHJJanA3SW1OeWRpSTZJa1ZrTWpVMU1Ua2lMQ0pyZEhraU9pSlBTMUFpTENKNElqb2lkMGRwVW5WdlRsVkZjbkZvU0d0WFoweGpiSFJHTW5sb1pYWjNhMGxTTTNaQk9FTkNTSGswTkhWU2F5SjlMQ0p3ZFhKd2IzTmxjeUk2V3lJalpIZHVMWE5wWnlKZGZTd2lkSGx3WlNJNklrUmxZMlZ1ZEhKaGJHbDZaV1JYWldKT2IyUmxJbjFkZlgxZExDSjFjR1JoZEdWRGIyMXRhWFJ0Wlc1MElqb2lSV2xDWTJneE4ySmZZbWN6TVhwelpIbDJkMVZtYlRreU1GWnVVbXBCU3pOa1JUaHZUR2RuYVRZd0xUaENRU0o5TENKemRXWm1hWGhFWVhSaElqcDdJbVJsYkhSaFNHRnphQ0k2SWtWcFF6RTFjbFJCVTBWWVNuaHFhMHBrU1Y5MGFtTlhhRTltWkhoU1RIZDVNSGg2ZDFsdE1XSmFVa2xTYldjaUxDSnlaV052ZG1WeWVVTnZiVzFwZEcxbGJuUWlPaUpGYVVOb016RnBVUzFaU1dJd1FYTkVibGRZYkVzM2NVZEVaWFJXVVU1aFMwMTJSM0F3Y1dwYVJtVlFjVEJuSW4xOSNkd24tc2lnIiwiYWxnIjoiRWREU0EifQ",
	"signature": "fqWuqWMEpdv2cf0KvozKasOmEd-QJLNlvvuMqNzVvJSbrrTx2EzKIr1IesfBoxdnmdeRkFivODsGVMuNVwjlBQ"
      }
    ]
  }
}
`

func NewTestDwn(t *testing.T) *Dwn {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "dwn-test")
	if err != nil {
		t.Fatal(err)
	}
	// Clean up after test
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	// Create subdirectories for blockstore and index
	blockstoreDir := filepath.Join(tempDir, "blockstore")
	indexDir := filepath.Join(tempDir, "index")

	err = os.MkdirAll(blockstoreDir, 0755)
	assert.NoError(t, err)

	err = os.MkdirAll(indexDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	didResolver := NewDidResolver([]DidMethodResolver{}, nil)

	dwn, err := NewDwn(DwnConfig{
		DidResolver:        didResolver,
		TenantGate:         NewAllowAllTenantGate(),
		MessageStore:       NewMemoryMessageStore(),
		DataStore:          NewMemoryDatastore(),
		EventLog:           NewMemoryEventLog(),
		BlockstoreLocation: blockstoreDir,
	})

	if err != nil {
		t.Fatal(err)
	}

	return dwn
}

func TestReadDid(t *testing.T) {
	content, err := os.ReadFile("testdata/did.json")
	assert.NoError(t, err)

	didInfo := fileDID{}
	err = json.Unmarshal(content, &didInfo)
	assert.NoError(t, err)

	assert.Equal(t, "ryan.rawson@abaxx.tech", didInfo.Id)
	assert.Equal(t, "ryan.rawson@abaxx.tech", didInfo.Metadata.Email)

	// find a key:
	// re-serialize a key to import it:
	printKey(t, didInfo, "#dwn-enc")
	printKey(t, didInfo, "#dwn-sig")

}

func printKey(t *testing.T, didInfo fileDID, keyid string) {
	keystruct, err := didInfo.GetKeyById(keyid)
	assert.NoError(t, err)

	t.Log("keystruct", keystruct)

	keyBytes, err := json.Marshal(keystruct)
	assert.NoError(t, err)

	k1, err := jwk.ParseKey(keyBytes)
	assert.NoError(t, err)
	t.Log("jwk k1", k1)

	t.Log("keyid", k1.KeyID())
	t.Log("keyops", k1.KeyOps())
	t.Log("keytype", k1.KeyType())
	t.Log("algo", k1.Algorithm().String())

	x, err := json.MarshalIndent(k1, "", "  ")

	assert.NoError(t, err)
	t.Log("serialize", string(x))
}

func GenerateTestCID(t *testing.T) string {
	// Create a random byte slice
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	assert.NoError(t, err, "Failed to generate random bytes")

	// Create a new CID using the random bytes
	mh, err := multihash.Sum(randomBytes, multihash.SHA2_256, -1)
	assert.NoError(t, err, "Failed to create multihash")

	c := cid.NewCidV1(cid.Raw, mh)
	return c.String()
}

func TestAuthSig(t *testing.T) {
	// Create a new DWN instance
	// dwn := NewTestDwn(t)

	// Generate a test CID
	testCid := GenerateTestCID(t)
	assert.NotEmpty(t, testCid, "Generated CID should not be empty")

	// Parse the auth signature
	var auth struct {
		Signature GeneralJws `json:"signature"`
	}
	err := json.Unmarshal([]byte(authSig), &auth)
	assert.NoError(t, err, "Failed to parse auth signature")
	assert.NotEmpty(t, auth.Signature.Payload, "Signature payload should not be empty")
	assert.NotEmpty(t, auth.Signature.Signatures, "Signature should have at least one signature")

}

func CreateTestIndexableKeyValues() IndexableKeyValues {
	indexableKV := IndexableKeyValues{
		"name":    S("John Doe"),
		"age":     I(30),
		"balance": F(1234.56),
		"active":  B(true),
	}

	// Using the indexableKV in a Put operation
	messageStore := NewMemoryMessageStore()
	tenant := Tenant("example-tenant")
	message := map[string]interface{}{
		"id":      "123456",
		"content": "Hello, World!",
	}

	err := messageStore.Put(tenant, message, indexableKV)
	if err != nil {
		fmt.Printf("Error putting message: %v", err)
	}

	return indexableKV
}

func TestCreateDIDAndSaveRecord(t *testing.T) {
	// Create a new DWN instance
	dwn := NewTestDwn(t)

	// Generate a new DID
	did, _, err := GenerateTestDID()
	assert.NoError(t, err, "Failed to generate test DID")
	assert.NotEmpty(t, did, "Generated DID should not be empty")

	// Create a test record
	record := &GenericMessage{
		descriptor: Descriptor{
			Method:      "CollectionsWrite",
			DataCid:     DataCid(GenerateTestCID(t)),
			DateCreated: time.Now().UTC().Format(time.RFC3339),
		},
		data: []byte("Test data for DWN record"),
	}

	// Create indexable key-values for the record
	indexableKeyValues := CreateTestIndexableKeyValues()

	// Save the record to the user's DWN
	err = dwn.messageStore.Put(Tenant(did), record, indexableKeyValues)
	assert.NoError(t, err, "Failed to save record to DWN")

	// Log the record before saving
	t.Logf("\nRecord to be saved: %+v\n\n", record)
	t.Logf("\nRecord CID: %+v\n\n", record.descriptor.DataCid)
	t.Logf("\nRecord descriptor: %+v\n\n", record.descriptor)
	t.Logf("\nRecord data: %s\n\n", string(record.data))
	t.Logf("\nIndexable key-values: %+v\n\n", indexableKeyValues)

	t.Logf("Type of Record: %T", record)

	// Retrieve the saved record
	retrievedRecord, err := dwn.messageStore.Get(Tenant(did), MessageCid(record.descriptor.DataCid))
	if err != nil {
		t.Fatalf("Error retrieving record: %v", err)
	}

	if retrievedRecord == nil {
		t.Fatalf("Retrieved record is nil")
	}

	genericMessage, ok := retrievedRecord.(*GenericMessage)
	if !ok {
		t.Fatalf("Retrieved record is not of type *GenericMessage, got: %T", retrievedRecord)
	}

	// Compare the retrieved record with the original
	assert.Equal(t, record.data, genericMessage.data, "Retrieved record data should match original")
	assert.Equal(t, record.descriptor.DataCid, genericMessage.descriptor.DataCid, "Retrieved record CID should match original")

	// Add more debugging information
	t.Logf("Original record: %+v", record)
	t.Logf("Retrieved record: %+v", genericMessage)
}

func GenerateTestDID() (string, *ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", nil, err
	}

	publicKeyBytes := elliptic.Marshal(privateKey.PublicKey.Curve, privateKey.PublicKey.X, privateKey.PublicKey.Y)
	didKey := fmt.Sprintf("did:key:z%s", base58.Encode(publicKeyBytes))

	return didKey, privateKey, nil
}
