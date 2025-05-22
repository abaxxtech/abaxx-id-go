package dwn

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlsoKnownAsExampleDID(t *testing.T) {
	// Load the example DID document
	content, err := os.ReadFile("testdata/did_with_also_known_as.json")
	require.NoError(t, err, "Failed to read example DID document")

	// Parse the DID document
	var exampleDID fileDID
	err = json.Unmarshal(content, &exampleDID)
	require.NoError(t, err, "Failed to parse example DID document")

	// Verify the DID document has the alsoKnownAs property
	alsoKnownAs := exampleDID.Metadata.Did.Document.AlsoKnownAs
	require.NotNil(t, alsoKnownAs, "alsoKnownAs property should exist")
	assert.Len(t, alsoKnownAs, 3, "alsoKnownAs should have 3 identifiers")

	// Check specific identifiers
	assert.Contains(t, alsoKnownAs, "https://example.com/users/exampleuser", "Should contain web URL")
	assert.Contains(t, alsoKnownAs, "mailto:user.example@abaxx.tech", "Should contain email URI")
	assert.Contains(t, alsoKnownAs, "did:key:z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK", "Should contain another DID")

	// Test validating the URIs in alsoKnownAs
	for _, uri := range alsoKnownAs {
		err := ValidateURI(uri)
		assert.NoError(t, err, "URI should be valid: %s", uri)
	}

	// Demonstrate bidirectional verification (conceptual)
	t.Run("Conceptual bidirectional verification", func(t *testing.T) {
		// In a real-world scenario, you would:

		// 1. For a web URL, you might check for a specific file on the website
		t.Log("For web URL: Fetch https://example.com/users/exampleuser/.well-known/did-configuration.json " +
			"and verify it contains the DID")

		// 2. For an email, you might send a verification email
		t.Log("For email: Send verification email to user.example@abaxx.tech " +
			"and wait for confirmation")

		// 3. For another DID, you would check its alsoKnownAs property
		t.Log("For other DID: Resolve did:key:z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK " +
			"and check its alsoKnownAs property for this DID")
	})

	// Demonstrate potential use cases
	t.Run("Use cases for alsoKnownAs", func(t *testing.T) {
		// Identity correlation - linking multiple identifiers for the same entity
		t.Log("Use case 1: Identity correlation - Showing that a DID is controlled by " +
			"the same entity that controls the email and website")

		// Legacy system integration - connecting DID to traditional identifiers
		t.Log("Use case 2: Legacy system integration - Allowing systems that understand " +
			"email but not DIDs to still identify the entity")

		// Cross-chain interoperability - linking DIDs across different methods/blockchains
		t.Log("Use case 3: Cross-chain interoperability - Establishing equivalence " +
			"between a did:ion and a did:key identifier")
	})
}
