package didion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/abaxxtech/abaxx-id-go/internal/crypto"
	"github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/did"
	"github.com/abaxxtech/abaxx-id-go/internal/dids/didcore"
	"github.com/abaxxtech/abaxx-id-go/internal/jwk"
)

const (
	defaultIONAPIEndpoint = "https://ion.msidentity.com/api/v1.0/"
)

type createOptions struct {
	keyManager  crypto.KeyManager
	algorithmID string
	apiEndpoint string
}

type CreateOption func(*createOptions)

type Resolver struct{}

func (r Resolver) Resolve(didURI string) (didcore.ResolutionResult, error) {
	return Resolve(didURI)
}

func (r Resolver) ResolveWithContext(ctx context.Context, didURI string) (didcore.ResolutionResult, error) {
	return ResolveWithContext(ctx, didURI)
}

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

	// Create ION DID request
	ionRequest, err := createIONRequest(document)
	if err != nil {
		return did.BearerDID{}, fmt.Errorf("failed to create ION request: %w", err)
	}

	// Send request to ION node
	ionDID, err := sendIONRequest(ctx, options.apiEndpoint, ionRequest)
	if err != nil {
		return did.BearerDID{}, fmt.Errorf("failed to send ION request: %w", err)
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
		Authentication: []string{"#key-1"},
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

type IONRequest struct {
	Type     string           `json:"type"`
	Document didcore.Document `json:"document"`
}

func createIONRequest(document didcore.Document) (IONRequest, error) {
	return IONRequest{
		Type:     "create",
		Document: document,
	}, nil
}

func sendIONRequest(ctx context.Context, apiEndpoint string, request IONRequest) (did.DID, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return did.DID{}, fmt.Errorf("failed to marshal ION request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint+"/operations", bytes.NewBuffer(jsonData))
	if err != nil {
		return did.DID{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return did.DID{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return did.DID{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		DID string `json:"did"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return did.DID{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return did.Parse(result.DID)
}

func Resolve(didURI string) (didcore.ResolutionResult, error) {
	return ResolveWithContext(context.Background(), didURI)
}

func ResolveWithContext(ctx context.Context, didURI string) (didcore.ResolutionResult, error) {
	// Validate the DID URI
	if !strings.HasPrefix(didURI, "did:ion:") {
		return didcore.ResolutionResult{}, fmt.Errorf("invalid ION DID: %s", didURI)
	}

	// Extract the ION-specific identifier
	// ionID := strings.TrimPrefix(didURI, "did:ion:")

	// Construct the resolution URL
	resolutionURL := fmt.Sprintf("%s/identifiers/%s", defaultIONAPIEndpoint, didURI)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", resolutionURL, nil)
	if err != nil {
		return didcore.ResolutionResult{}, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Send the request
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return didcore.ResolutionResult{}, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		return didcore.ResolutionResult{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var result didcore.ResolutionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return didcore.ResolutionResult{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// Update updates an existing ION DID
func Update(didURI string, document didcore.Document, keyManager crypto.KeyManager) error {
	return UpdateWithContext(context.Background(), didURI, document, keyManager)
}

func UpdateWithContext(ctx context.Context, didURI string, document didcore.Document, keyManager crypto.KeyManager) error {
	updateRequest := IONRequest{
		Type:     "update",
		Document: document,
	}

	// Sign the update request
	// Note: You'll need to implement the signing logic based on ION's requirements

	jsonData, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal update request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", defaultIONAPIEndpoint+"operations", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// Deactivate deactivates an existing ION DID
func Deactivate(didURI string, keyManager crypto.KeyManager) error {
	return DeactivateWithContext(context.Background(), didURI, keyManager)
}

func DeactivateWithContext(ctx context.Context, didURI string, keyManager crypto.KeyManager) error {
	deactivateRequest := IONRequest{
		Type: "deactivate",
	}

	// Sign the deactivate request
	// Note: You'll need to implement the signing logic based on ION's requirements

	jsonData, err := json.Marshal(deactivateRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal deactivate request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", defaultIONAPIEndpoint+"operations", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
