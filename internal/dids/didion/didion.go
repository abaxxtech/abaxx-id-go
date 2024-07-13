package didion

import (
	"context"
	"fmt"

	"github.com/abaxxtech/abaxx-id-go/internal/crypto"
	"github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/did"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didcore"
	"github.com/abaxxtech/abaxx-id-go/internal/jwk"
)

const (
	defaultIONAPIEndpoint = "https://ion.abaxx.id" // @todo
)

type createOptions struct {
	keyManager  crypto.KeyManager
	algorithmID string
	apiEndpoint string
}

type CreateOption func(*createOptions)

func KeyManager(km crypto.KeyManager) CreateOption {
	return func(o *createOptions) {
		o.keyManager = km
	}
}

func AlgorithmID(alg string) CreateOption {
	return func(o *createOptions) {
		o.algorithmID = alg
	}
}

func APIEndpoint(endpoint string) CreateOption {
	return func(o *createOptions) {
		o.apiEndpoint = endpoint
	}
}

func Create(opts ...CreateOption) (did.BearerDID, error) {
	return CreateWithContext(context.Background(), opts...)
}

func CreateWithContext(ctx context.Context, opts ...CreateOption) (did.BearerDID, error) {
	options := createOptions{
		keyManager:  crypto.NewLocalKeyManager(),
		algorithmID: dsa.AlgorithmIDED25519,
		apiEndpoint: defaultIONAPIEndpoint,
	}

	for _, opt := range opts {
		opt(&options)
	}

	// Generate key pair
	keyID, err := options.keyManager.GeneratePrivateKey(options.algorithmID)
	if err != nil {
		return did.BearerDID{}, fmt.Errorf("failed to generate private key: %w", err)
	}

	publicKey, err := options.keyManager.GetPublicKey(keyID)
	if err != nil {
		return did.BearerDID{}, fmt.Errorf("failed to get public key: %w", err)
	}

	// Create ION DID document
	document := createIONDocument(publicKey)

	// For now, we'll return a placeholder DID
	ionDID := did.DID{
		Method: "ion",
		URI:    "did:ion:placeholder",
	}

	return did.BearerDID{
		DID:        ionDID,
		KeyManager: options.keyManager,
		Document:   document,
	}, nil
}

func createIONDocument(publicKey jwk.JWK) didcore.Document {
	document := didcore.Document{
		Context: []string{"https://www.w3.org/ns/did/v1"},
		ID:      "did:ion:placeholder",
		VerificationMethod: []didcore.VerificationMethod{
			{
				ID:           "#key-1",
				Type:         "JsonWebKey2020",
				Controller:   "did:ion:placeholder",
				PublicKeyJwk: &publicKey,
			},
		},
		Service: []didcore.Service{
			{
				ID:              "#dwn",
				Type:            "DecentralizedWebNode",
				ServiceEndpoint: []string{"http://localhost:8085"},
			},
		},
	}
	return document
}

// TODO: Implement ION-specific operations
// 1. Create ION DID request
// 2. Resolve ION DID
// 3. Wait for anchoring and retrieve long-form DID

// ionDIDRequest, err := createIONDIDRequest(document)
// if err != nil {
// 	return did.BearerDID{}, fmt.Errorf("failed to create ION DID request: %w", err)
// }

// resolvedDID, err := resolveIONDID(ionDIDRequest)
// if err != nil {
// 	return did.BearerDID{}, fmt.Errorf("failed to resolve ION DID: %w", err)
// }

// longFormDID, err := waitForAnchoringAndRetrieveLongFormDID(resolvedDID)
// if err != nil {
// 	return did.BearerDID{}, fmt.Errorf("failed to retrieve long-form DID: %w", err)
// }
