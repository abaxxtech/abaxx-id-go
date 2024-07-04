package did

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Cache interface
type Cache interface {
	Get(key string) (DidResolutionResult, bool)
	Set(key string, value DidResolutionResult)
}

// DidMethodResolver interface
type DidMethodResolver interface {
	Method() string
	Resolve(did string) (DidResolutionResult, error)
}

// DidResolutionResult struct
type DidResolutionResult struct {
	DidDocument           interface{}
	DidResolutionMetadata struct {
		Error string
	}
}

// MemoryCache struct
type MemoryCache struct {
	data      map[string]DidResolutionResult
	expiry    time.Duration
	timestamp map[string]time.Time
}

func NewMemoryCache(expiry time.Duration) *MemoryCache {
	return &MemoryCache{
		data:      make(map[string]DidResolutionResult),
		expiry:    expiry,
		timestamp: make(map[string]time.Time),
	}
}

func (c *MemoryCache) Get(key string) (DidResolutionResult, bool) {
	if val, found := c.data[key]; found {
		if time.Since(c.timestamp[key]) < c.expiry {
			return val, true
		}
		delete(c.data, key)
		delete(c.timestamp, key)
	}
	return DidResolutionResult{}, false
}

func (c *MemoryCache) Set(key string, value DidResolutionResult) {
	c.data[key] = value
	c.timestamp[key] = time.Now()
}

// DidResolver struct
type DidResolver struct {
	didResolvers map[string]DidMethodResolver
	cache        Cache
}

func NewDidResolver(resolvers []DidMethodResolver, cache Cache) *DidResolver {
	if cache == nil {
		cache = NewMemoryCache(10 * time.Minute)
	}

	if resolvers == nil || len(resolvers) == 0 {
		resolvers = []DidMethodResolver{
			// NewDidIonResolver(),
			// NewDidKeyResolver(),
		}
	}

	didResolvers := make(map[string]DidMethodResolver)
	for _, resolver := range resolvers {
		didResolvers[resolver.Method()] = resolver
	}

	return &DidResolver{
		didResolvers: didResolvers,
		cache:        cache,
	}
}

// TODO implement/remove this
func Validate(did string) error {
	return nil
}

func extractMethod(did string) string {
	return ""
}


func (r *DidResolver) Resolve(did string) (DidResolutionResult, error) {
	if err := Validate(did); err != nil {
		return DidResolutionResult{}, err
	}

	if result, found := r.cache.Get(did); found {
		return result, nil
	}

	method := extractMethod(did)
	resolver, exists := r.didResolvers[method]
	if !exists {
		return DidResolutionResult{}, fmt.Errorf("no resolver found for method: %s", method)
	}

	result, err := resolver.Resolve(did)
	if err != nil {
		return DidResolutionResult{}, err
	}

	r.cache.Set(did, result)
	return result, nil
}

// ParseDidDocument parses a DID document from a JSON file
func ParseDidDocument(filePath string) (DidResolutionResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return DidResolutionResult{}, err
	}

	var didDoc DidResolutionResult
	if err := json.Unmarshal(data, &didDoc); err != nil {
		return DidResolutionResult{}, err
	}

	return didDoc, nil
}

// FetchDidDocument fetches a DID document from a URL
func FetchDidDocument(url string) (DidResolutionResult, error) {
	resp, err := http.Get(url)
	if err != nil {
		return DidResolutionResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DidResolutionResult{}, errors.New("failed to fetch DID document")
	}

	var didDoc DidResolutionResult
	if err := json.NewDecoder(resp.Body).Decode(&didDoc); err != nil {
		return DidResolutionResult{}, err
	}

	return didDoc, nil
}

func (r *DidResolver) ResolveFromFile(filePath string) (DidResolutionResult, error) {
	didDoc, err := ParseDidDocument(filePath)
	if err != nil {
		return DidResolutionResult{}, err
	}
	r.cache.Set(didDoc.DidDocument.(string), didDoc)
	return didDoc, nil
}

func (r *DidResolver) ResolveFromURL(url string) (DidResolutionResult, error) {
	didDoc, err := FetchDidDocument(url)
	if err != nil {
		return DidResolutionResult{}, err
	}
	r.cache.Set(didDoc.DidDocument.(string), didDoc)
	return didDoc, nil
}
