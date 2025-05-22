package dwn

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// Errors related to alsoKnownAs operations
var (
	ErrInvalidURI         = errors.New("invalid URI format")
	ErrIdentifierExists   = errors.New("identifier already exists in alsoKnownAs")
	ErrIdentifierNotFound = errors.New("identifier not found in alsoKnownAs")
	ErrCannotResolveDoc   = errors.New("cannot resolve DID document")
	ErrInvalidAlsoKnownAs = errors.New("invalid alsoKnownAs property")
)

// ValidateURI checks if the provided string is a valid URI
func ValidateURI(uri string) error {
	// Check if it's a DID
	if strings.HasPrefix(uri, "did:") {
		parts := strings.Split(uri, ":")
		if len(parts) < 3 {
			return ErrInvalidURI
		}
		return nil
	}

	// Otherwise validate as regular URI
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidURI, err.Error())
	}
	return nil
}

// AddAlsoKnownAs adds a new identifier to the alsoKnownAs property of a DID document
func (d *Dwn) AddAlsoKnownAs(did DID, alternativeID string) error {
	// Validate the URI format
	if err := ValidateURI(alternativeID); err != nil {
		return err
	}

	// Resolve the DID document
	didResult, err := d.didResolver.Resolve(string(did))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCannotResolveDoc, err.Error())
	}

	// Get the DID document
	doc, ok := didResult.DidDocument.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse DID document: unexpected document type %T", didResult.DidDocument)
	}

	// Check if alsoKnownAs already exists
	var alsoKnownAs []string
	if existingIDs, exists := doc["alsoKnownAs"]; exists {
		if ids, ok := existingIDs.([]string); ok {
			alsoKnownAs = ids
			// Check if the identifier already exists
			for _, id := range alsoKnownAs {
				if id == alternativeID {
					return ErrIdentifierExists
				}
			}
		} else if ids, ok := existingIDs.([]interface{}); ok {
			// Handle case where it might be []interface{} containing strings
			for _, id := range ids {
				if strID, ok := id.(string); ok {
					alsoKnownAs = append(alsoKnownAs, strID)
					if strID == alternativeID {
						return ErrIdentifierExists
					}
				}
			}
		} else {
			return ErrInvalidAlsoKnownAs
		}
	}

	// Add the new identifier
	alsoKnownAs = append(alsoKnownAs, alternativeID)
	doc["alsoKnownAs"] = alsoKnownAs

	// Update the DID document in the appropriate way (this would depend on your DID method)
	// This is a placeholder that would need to be implemented based on your specific DID method
	// For example, this might involve creating a transaction for a blockchain-based DID

	return nil
}

// RemoveAlsoKnownAs removes an identifier from the alsoKnownAs property of a DID document
func (d *Dwn) RemoveAlsoKnownAs(did DID, identifierToRemove string) error {
	// Resolve the DID document
	didResult, err := d.didResolver.Resolve(string(did))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCannotResolveDoc, err.Error())
	}

	// Get the DID document
	doc, ok := didResult.DidDocument.(map[string]interface{})
	if !ok {
		return fmt.Errorf("failed to parse DID document: unexpected document type %T", didResult.DidDocument)
	}

	// Check if alsoKnownAs exists
	existingIDs, exists := doc["alsoKnownAs"]
	if !exists {
		return ErrIdentifierNotFound
	}

	// Extract the current identifiers
	var alsoKnownAs []string
	var found bool

	if ids, ok := existingIDs.([]string); ok {
		for _, id := range ids {
			if id == identifierToRemove {
				found = true
			} else {
				alsoKnownAs = append(alsoKnownAs, id)
			}
		}
	} else if ids, ok := existingIDs.([]interface{}); ok {
		for _, id := range ids {
			if strID, ok := id.(string); ok {
				if strID == identifierToRemove {
					found = true
				} else {
					alsoKnownAs = append(alsoKnownAs, strID)
				}
			}
		}
	} else {
		return ErrInvalidAlsoKnownAs
	}

	if !found {
		return ErrIdentifierNotFound
	}

	// Update the alsoKnownAs property
	if len(alsoKnownAs) > 0 {
		doc["alsoKnownAs"] = alsoKnownAs
	} else {
		delete(doc, "alsoKnownAs")
	}

	// Update the DID document in the appropriate way (this would depend on your DID method)
	// This is a placeholder that would need to be implemented based on your specific DID method

	return nil
}

// GetAlsoKnownAs retrieves all alternative identifiers for a DID
func (d *Dwn) GetAlsoKnownAs(did DID) ([]string, error) {
	// Resolve the DID document
	didResult, err := d.didResolver.Resolve(string(did))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCannotResolveDoc, err.Error())
	}

	// Get the DID document
	doc, ok := didResult.DidDocument.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse DID document: unexpected document type %T", didResult.DidDocument)
	}

	// Extract alsoKnownAs
	var alsoKnownAs []string
	if existingIDs, exists := doc["alsoKnownAs"]; exists {
		if ids, ok := existingIDs.([]string); ok {
			alsoKnownAs = ids
		} else if ids, ok := existingIDs.([]interface{}); ok {
			for _, id := range ids {
				if strID, ok := id.(string); ok {
					alsoKnownAs = append(alsoKnownAs, strID)
				}
			}
		}
	}

	return alsoKnownAs, nil
}

// VerifyBidirectionalAlsoKnownAs verifies that the linked identifier also references this DID
// This helps establish trust in the alsoKnownAs relationship
func (d *Dwn) VerifyBidirectionalAlsoKnownAs(did DID, alternativeID string) (bool, error) {
	// First check if the alternativeID is in the DID's alsoKnownAs
	identifiers, err := d.GetAlsoKnownAs(did)
	if err != nil {
		return false, err
	}

	foundInDID := false
	for _, id := range identifiers {
		if id == alternativeID {
			foundInDID = true
			break
		}
	}

	if !foundInDID {
		return false, nil
	}

	// If alternativeID is a DID, check if it references back
	if strings.HasPrefix(alternativeID, "did:") {
		altDIDResult, err := d.didResolver.Resolve(alternativeID)
		if err != nil {
			return false, err
		}

		altDoc, ok := altDIDResult.DidDocument.(map[string]interface{})
		if !ok {
			return false, fmt.Errorf("failed to parse alternative DID document")
		}

		if altAlsoKnownAs, exists := altDoc["alsoKnownAs"]; exists {
			if ids, ok := altAlsoKnownAs.([]string); ok {
				for _, id := range ids {
					if id == string(did) {
						return true, nil
					}
				}
			} else if ids, ok := altAlsoKnownAs.([]interface{}); ok {
				for _, id := range ids {
					if strID, ok := id.(string); ok && strID == string(did) {
						return true, nil
					}
				}
			}
		}
		return false, nil
	}

	// For non-DID URIs, verification would depend on the URI type
	// For example, for a website, you might check for a specific file or metadata
	// This is a placeholder for such implementation
	return false, fmt.Errorf("verification of non-DID URIs not implemented")
}
