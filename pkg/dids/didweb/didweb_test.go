package didweb_test

import (
	"testing"

	"github.com/abaxxtech/abaxx-id-go/pkg/dids/didcore"
	"github.com/abaxxtech/abaxx-id-go/pkg/dids/didweb"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	bearerDID, err := didweb.Create("localhost:8080")
	assert.NoError(t, err)

	assert.NotEqual(t, didcore.Document{}, bearerDID.Document)

	document := bearerDID.Document
	assert.Equal(t, "did:web:localhost%3A8080", document.ID)
	assert.Equal(t, 1, len(document.VerificationMethod))
}

func TestTransformID(t *testing.T) {
	var vectors = []struct {
		input  string
		output string
		err    bool
	}{
		{
			input:  "example.com:user:alice",
			output: "https://example.com/user/alice/did.json",
			err:    false,
		},
		{
			input:  "localhost%3A8080:user:alice",
			output: "http://localhost:8080/user/alice/did.json",
			err:    false,
		},
		{
			input:  "192.168.1.100%3A8892:ingress",
			output: "http://192.168.1.100:8892/ingress/did.json",
			err:    false,
		},
		{
			input:  "www.linkedin.com",
			output: "https://www.linkedin.com/.well-known/did.json",
			err:    false,
		},
	}

	for _, v := range vectors {
		t.Run(v.input, func(t *testing.T) {
			output, err := didweb.TransformID(v.input)
			assert.NoError(t, err)
			assert.Equal(t, v.output, output)
		})
	}
}
