# The `alsoKnownAs` Property in DIDs

## Overview

The `alsoKnownAs` property is a standard feature in the W3C DID Core specification (https://www.w3.org/TR/did-1.0/#also-known-as) that allows a DID subject to establish connections with other identifiers representing the same entity. This document explains how to use the `alsoKnownAs` property in the Abaxx DWN (Decentralized Web Node) implementation.

## Purpose

The `alsoKnownAs` property serves several important purposes:

1. **Identity Correlation**: Establish connections between different identifiers that represent the same entity.
2. **Legacy System Integration**: Connect DIDs to traditional identifiers like email addresses, domain names, etc.
3. **Cross-chain Interoperability**: Link DIDs across different methods or blockchains.
4. **Identifier Continuity**: Maintain identifier continuity when migrating from one DID method to another.

## Format

The `alsoKnownAs` property is an array of URIs that can include:

- Other DIDs (e.g., `did:ion:123...`, `did:key:456...`)
- Web URLs (e.g., `https://example.com/users/alice`)
- Email URIs (e.g., `mailto:alice@example.com`)
- Other URI formats that identify the same entity

## Example DID Document with `alsoKnownAs`

```json
{
  "id": "did:ion:EiAFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Avw",
  "alsoKnownAs": [
    "https://example.com/users/exampleuser",
    "mailto:user.example@abaxx.tech",
    "did:key:z6MkhaXgBZDvotDkL5257faiztiGiC2QtKLGpbnnEGta2doK"
  ],
  "verificationMethod": [
    // verification methods
  ],
  "service": [
    // services
  ]
}
```

## API Usage

### Adding an Alternative Identifier

```go
// Add a web URL as an alternative identifier
err := dwnInstance.AddAlsoKnownAs("did:ion:123...", "https://example.com/users/alice")
if err != nil {
    // Handle error
}

// Add another DID as an alternative identifier
err = dwnInstance.AddAlsoKnownAs("did:ion:123...", "did:key:456...")
if err != nil {
    // Handle error
}

// Add an email as an alternative identifier
err = dwnInstance.AddAlsoKnownAs("did:ion:123...", "mailto:alice@example.com")
if err != nil {
    // Handle error
}
```

### Getting Alternative Identifiers

```go
// Get all alternative identifiers for a DID
identifiers, err := dwnInstance.GetAlsoKnownAs("did:ion:123...")
if err != nil {
    // Handle error
}

for _, id := range identifiers {
    fmt.Printf("Alternative identifier: %s\n", id)
}
```

### Removing an Alternative Identifier

```go
// Remove an alternative identifier
err := dwnInstance.RemoveAlsoKnownAs("did:ion:123...", "https://example.com/users/alice")
if err != nil {
    // Handle error
}
```

### Verifying Bidirectional Relationships

For maximum trust, it's recommended to verify that alternative identifiers also reference back to the DID:

```go
// Check if the relationship is bidirectional
isBidirectional, err := dwnInstance.VerifyBidirectionalAlsoKnownAs("did:ion:123...", "did:key:456...")
if err != nil {
    // Handle error
}

if isBidirectional {
    fmt.Println("The relationship is verified in both directions")
} else {
    fmt.Println("The alternative identifier does not reference back to this DID")
}
```

## Best Practices

1. **Verify Ownership**: Before adding an alternative identifier, verify that the DID controller actually controls that identifier.
2. **Use Bidirectional References**: When possible, establish bidirectional references between identifiers.
3. **Limit Personal Information**: Be cautious about including identifiers that reveal personal information.
4. **Regular Validation**: Periodically validate that alternative identifiers are still valid and controlled by the same entity.
5. **Consider Privacy**: Be aware that linking identifiers may reduce privacy by making correlation easier.

## Use Cases

### Identity Migration

When migrating from one DID method to another, you can use `alsoKnownAs` to indicate that both DIDs represent the same entity:

```json
{
  "id": "did:ion:newIdentifier",
  "alsoKnownAs": ["did:method:oldIdentifier"]
}
```

### Website Authentication

Link a DID to a domain for authentication purposes:

```json
{
  "id": "did:ion:123...",
  "alsoKnownAs": ["https://example.com"]
}
```

The website can verify this relationship by hosting a DID configuration document at a well-known location (e.g., `/.well-known/did-configuration.json`).

### Cross-platform Identity

Link identities across different platforms to create a comprehensive digital identity:

```json
{
  "id": "did:ion:123...",
  "alsoKnownAs": [
    "https://twitter.com/username",
    "https://github.com/username",
    "https://linkedin.com/in/username"
  ]
}
```

## Security Considerations

1. **Spoofing Risk**: Without proper verification, malicious actors could claim relationships with identifiers they don't control.
2. **Privacy Leakage**: Linking identifiers can lead to correlation of previously separate identities.
3. **Verification Challenges**: Different identifier types require different verification mechanisms.

## Related Documentation

- [W3C DID Core Specification: alsoKnownAs](https://www.w3.org/TR/did-1.0/#also-known-as)
- [DID Configuration Specification](https://identity.foundation/.well-known/resources/did-configuration/) 