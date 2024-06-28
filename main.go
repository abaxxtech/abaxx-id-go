package main

import (
	"fmt"
	"time"

	"github.com/abaxxtech/abaxx-id-go/internal/did"
)

func main() {
	cache := did.NewMemoryCache(10 * time.Minute)
	resolver := did.NewDidResolver(nil, cache)

	// Resolve DID from file
	didDoc, err := resolver.ResolveFromFile("path/to/did-document.json")
	if err != nil {
		fmt.Println("Error resolving DID from file:", err)
		return
	}
	fmt.Println("Resolved DID from file:", didDoc)

	// Resolve DID from URL
	didDoc, err = resolver.ResolveFromURL("file://json-schemas/did-document.json")
	if err != nil {
		fmt.Println("Error resolving DID from URL:", err)
		return
	}
	fmt.Println("Resolved DID from URL:", didDoc)
}
