package dwn

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURI(t *testing.T) {
	testCases := []struct {
		name    string
		uri     string
		isValid bool
	}{
		{"Valid DID", "did:ion:1234567890", true},
		{"Valid HTTP URL", "https://example.com", true},
		{"Invalid DID - missing method", "did:", false},
		{"Invalid URL", "not-a-url", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateURI(tc.uri)
			if tc.isValid {
				assert.NoError(t, err, "Expected valid URI")
			} else {
				assert.Error(t, err, "Expected invalid URI")
			}
		})
	}
}

func TestAlsoKnownAs(t *testing.T) {
	// Setup a test DID document with mock resolver
	mockDoc := map[string]interface{}{
		"id":             "did:example:123",
		"authentication": []string{"did:example:123#keys-1"},
		"service": []map[string]interface{}{
			{
				"id":              "#dwn",
				"type":            "DecentralizedWebNode",
				"serviceEndpoint": map[string]interface{}{"nodes": []string{"https://example.com/dwn"}},
			},
		},
	}

	// Create a mock resolver that returns our test document
	mockResolver := &MockDidResolver{
		didDoc: mockDoc,
	}

	// Create DWN instance with mock resolver
	dwn := &Dwn{
		didResolver: &DidResolver{
			didResolvers: map[string]DidMethodResolver{
				"example": mockResolver,
			},
			cache: NewMemoryCache(0), // No cache expiry for testing
		},
	}

	// Test adding alsoKnownAs
	t.Run("AddAlsoKnownAs", func(t *testing.T) {
		// Test adding a valid DID
		err := dwn.AddAlsoKnownAs("did:example:123", "did:example:456")
		assert.NoError(t, err, "Failed to add valid DID")

		// Verify it was added
		alsoKnownAs, ok := mockDoc["alsoKnownAs"].([]string)
		assert.True(t, ok, "alsoKnownAs should be []string")
		assert.Contains(t, alsoKnownAs, "did:example:456", "alsoKnownAs should contain added DID")

		// Test adding a valid URL
		err = dwn.AddAlsoKnownAs("did:example:123", "https://example.com/profile")
		assert.NoError(t, err, "Failed to add valid URL")

		// Verify it was added
		alsoKnownAs, ok = mockDoc["alsoKnownAs"].([]string)
		assert.True(t, ok, "alsoKnownAs should be []string")
		assert.Contains(t, alsoKnownAs, "https://example.com/profile", "alsoKnownAs should contain added URL")

		// Test adding a duplicate
		err = dwn.AddAlsoKnownAs("did:example:123", "did:example:456")
		assert.Error(t, err, "Should fail when adding duplicate")
		assert.ErrorIs(t, err, ErrIdentifierExists)

		// Test adding an invalid URI
		err = dwn.AddAlsoKnownAs("did:example:123", "not-a-valid-uri")
		assert.Error(t, err, "Should fail with invalid URI")
	})

	// Test getting alsoKnownAs
	t.Run("GetAlsoKnownAs", func(t *testing.T) {
		identifiers, err := dwn.GetAlsoKnownAs("did:example:123")
		assert.NoError(t, err, "Failed to get alsoKnownAs")
		assert.Len(t, identifiers, 2, "Should have 2 identifiers")
		assert.Contains(t, identifiers, "did:example:456", "Should contain the first added identifier")
		assert.Contains(t, identifiers, "https://example.com/profile", "Should contain the second added identifier")
	})

	// Test removing alsoKnownAs
	t.Run("RemoveAlsoKnownAs", func(t *testing.T) {
		// Remove one identifier
		err := dwn.RemoveAlsoKnownAs("did:example:123", "did:example:456")
		assert.NoError(t, err, "Failed to remove alsoKnownAs")

		// Verify it was removed
		identifiers, err := dwn.GetAlsoKnownAs("did:example:123")
		assert.NoError(t, err, "Failed to get alsoKnownAs after removal")
		assert.Len(t, identifiers, 1, "Should have 1 identifier left")
		assert.NotContains(t, identifiers, "did:example:456", "Should not contain removed identifier")
		assert.Contains(t, identifiers, "https://example.com/profile", "Should contain the remaining identifier")

		// Try to remove a non-existent identifier
		err = dwn.RemoveAlsoKnownAs("did:example:123", "did:example:789")
		assert.Error(t, err, "Should fail when removing non-existent identifier")
		assert.ErrorIs(t, err, ErrIdentifierNotFound)

		// Remove the last identifier
		err = dwn.RemoveAlsoKnownAs("did:example:123", "https://example.com/profile")
		assert.NoError(t, err, "Failed to remove last alsoKnownAs")

		// Verify it was removed
		identifiers, err = dwn.GetAlsoKnownAs("did:example:123")
		assert.NoError(t, err, "Failed to get alsoKnownAs after removal")
		assert.Empty(t, identifiers, "Should have no identifiers left")
		_, exists := mockDoc["alsoKnownAs"]
		assert.False(t, exists, "alsoKnownAs property should be removed when empty")
	})
}

// MockDidResolver is a simple resolver that returns a predefined document
type MockDidResolver struct {
	didDoc map[string]interface{}
}

func (r *MockDidResolver) Method() string {
	return "example"
}

func (r *MockDidResolver) Resolve(did string) (DidResolutionResult, error) {
	return DidResolutionResult{
		DidDocument: r.didDoc,
	}, nil
}

// Helper function to pretty-print a DID document
func prettyPrintDocument(t *testing.T, doc map[string]interface{}) {
	bytes, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		t.Logf("Failed to marshal document: %v", err)
		return
	}
	t.Logf("DID Document:\n%s", string(bytes))
}
