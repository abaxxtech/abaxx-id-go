package diddht

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"io"

	"github.com/abaxxtech/abaxx-id-go/internal/crypto"
	"github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didcore"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/dns/dnsmessage"
)

type DHTTXTResourceOpt func() dnsmessage.Resource

func WithDNSRecord(name, body string) DHTTXTResourceOpt {
	return func() dnsmessage.Resource {
		return dnsmessage.Resource{
			Header: dnsmessage.ResourceHeader{
				Name: dnsmessage.MustNewName(name),
				Type: dnsmessage.TypeTXT,
				TTL:  7200,
			},
			Body: &dnsmessage.TXTResource{
				TXT: []string{
					body,
				},
			},
		}
	}
}

func TestCreate(t *testing.T) {
	tests := map[string]struct {
		didURI         string
		expectedResult string
		didDocData     string
		keys           []verificationMethodOption
	}{
		"": {
			didURI:         "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko",
			expectedResult: "1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko",
			keys: []verificationMethodOption{
				{
					algorithmID: dsa.AlgorithmIDED25519,
					purposes:    []didcore.Purpose{didcore.PurposeAssertion, didcore.PurposeAuthentication, didcore.PurposeCapabilityDelegation, didcore.PurposeCapabilityInvocation},
				},
			},
			didDocData: `{
				"id": "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko",
				"verificationMethod": [
				  {
					"id": "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko#0",
					"type": "JsonWebKey",
					"controller": "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko",
					"publicKeyJwk": {
					  "kty": "OKP",
					  "crv": "Ed25519",
					  "x": "lSuMYhg12IMawqFut-2URA212Nqe8-WEB7OBlam5oBU",
					  "kid": "2Jr7faCpoEgHvy5HXH32z-MH_0CRToO9NllZtemVvNo"
					}
				  }
				],
				"service": [
				  {
					"id": "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko#dwn",
					"type": "DecentralizedWebNode",
					"serviceEndpoint": ["https://example.com/dwn1", "https://example.com/dwn2"]
				  }
				],
				"authentication": [
				  "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko#0"
				],
				"assertionMethod": [
				  "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko#0"
				],
				"capabilityDelegation": [
				  "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko#0"
				],
				"capabilityInvocation": [
				  "did:dht:1wiaaaoagzceggsnwfzmx5cweog5msg4u536mby8sqy3mkp3wyko#0"
				]
			  }`,
		},
	}

	// setting up a fake relay that stores did documents on publish, and responds with the bencoded did document on resolve
	mockedRes := map[string][]byte{}
	relay := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		did := "did:dht:" + r.URL.Path[1:]
		defer r.Body.Close()

		// create branch
		if r.Method != http.MethodGet {
			packagedDid, err := io.ReadAll(r.Body)
			assert.NoError(t, err)
			mockedRes[did] = packagedDid
			w.WriteHeader(http.StatusOK)
			return
		}

		// resolve branch
		expectedBuf, ok := mockedRes[did]
		if !ok {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		_, err := w.Write(expectedBuf)
		assert.NoError(t, err)
	}))

	defer relay.Close()
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var didDoc didcore.Document
			assert.NoError(t, json.Unmarshal([]byte(test.didDocData), &didDoc))
			keyMgr := crypto.NewLocalKeyManager()

			var opts []CreateOption
			opts = []CreateOption{Gateway(relay.URL, http.DefaultClient), KeyManager(keyMgr)}
			for _, service := range didDoc.Service {
				opts = append(opts, Service(service.ID, service.Type, service.ServiceEndpoint...))
			}
			for _, key := range test.keys {
				opts = append(opts, PrivateKey(key.algorithmID, key.purposes...))
			}

			createdDid, err := Create(opts...)
			assert.NoError(t, err)
			assert.NotZero(t, createdDid.KeyManager)
			resolver := NewResolver(relay.URL, http.DefaultClient)
			result, err := resolver.Resolve(createdDid.URI)
			assert.Equal(t, len(createdDid.Document.VerificationMethod), 2)
			assert.NoError(t, err)
			assert.Equal(t, len(createdDid.Document.Authentication), 2)
			assert.Equal(t, len(createdDid.Document.AssertionMethod), 2)
			assert.Equal(t, len(createdDid.Document.CapabilityDelegation), 2)
			assert.Equal(t, len(createdDid.Document.CapabilityInvocation), 2)
			assert.Equal(t, createdDid.Document.Service, result.Document.Service)
		})
	}
}
