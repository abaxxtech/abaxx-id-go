package didjwk_test

import (
	"testing"

	"github.com/abaxxtech/abaxx-id-go/pkg/dids/didcore"
	"github.com/abaxxtech/abaxx-id-go/pkg/dids/didjwk"
	"github.com/abaxxtech/abaxx-id-go/pkg/jwk"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	did, err := didjwk.Create()
	assert.NoError(t, err)

	assert.Equal(t, "jwk", did.Method)
	assert.True(t, did.Fragment == "", "expected fragment to be empty")
	assert.True(t, did.Path == "", "expected path to be empty")
	assert.True(t, did.Query == "", "expected query to be empty")
	assert.True(t, did.ID != "", "expected id to be non-empty")
	assert.Equal(t, "did:jwk:"+did.ID, did.URI)
}

func TestResolveDIDJWK(t *testing.T) {
	resolver := &didjwk.Resolver{}
	result, err := resolver.Resolve("did:jwk:eyJraWQiOiJ1cm46aWV0ZjpwYXJhbXM6b2F1dGg6andrLXRodW1icHJpbnQ6c2hhLTI1NjpGZk1iek9qTW1RNGVmVDZrdndUSUpqZWxUcWpsMHhqRUlXUTJxb2JzUk1NIiwia3R5IjoiT0tQIiwiY3J2IjoiRWQyNTUxOSIsImFsZyI6IkVkRFNBIiwieCI6IkFOUmpIX3p4Y0tCeHNqUlBVdHpSYnA3RlNWTEtKWFE5QVBYOU1QMWo3azQifQ")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(result.Document.VerificationMethod))

	vm := result.Document.VerificationMethod[0]
	assert.True(t, vm != didcore.VerificationMethod{}, "expected verification method to be non-empty")
	assert.NotEqual(t, *vm.PublicKeyJwk, jwk.JWK{}, "expected publicKeyJwk to be non-empty")

	assert.Equal(t, 1, len(result.Document.Authentication))
	assert.Equal(t, 1, len(result.Document.AssertionMethod))
}
