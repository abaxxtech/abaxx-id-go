# Digital Title Legal Agreements as Verifiable Credential Claims

This document explains how to use digital title legal agreements as claims in verifiable credentials within the Abaxx ID protocol.

## Overview

The Abaxx ID protocol supports multiple approaches for representing digital title legal agreements as verifiable credential claims:

1. **Full Digital Title as Credential Subject** - Embed the entire legal agreement
2. **Reference-Based Claims** - Create lightweight credentials that reference the full agreement  
3. **Generic Claims** - Use flexible key-value pairs for varying structures

## Approach 1: Full Digital Title as Credential Subject

This approach embeds the complete digital title legal agreement as the credential subject using strongly-typed Go structs.

### Benefits
- **Complete Data**: All agreement details are included in the credential
- **Type Safety**: Strong typing ensures data integrity
- **Self-Contained**: No need to resolve external references

### Usage

```go
import (
    "github.com/abaxxtech/abaxx-id-go/pkg/vc"
    "github.com/abaxxtech/abaxx-id-go/pkg/dids/didjwk"
)

// Create the digital title claims
digitalTitle := &vc.DigitalTitleClaims{
    ID:        "did:dht:EiCFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Xco",
    TitleID:   "770f9510-e29b-41d4-a716-446655440002",
    TitleType: "legal-agreement",
    Asset: vc.Asset{
        Identifier:  "TRUST-INDENTURE-XCO-ABAXX-2024-001",
        Name:        "Gold Collateral Trust Indenture and Margin Lending Agreement",
        Description: "Trust arrangement where Xco pledges gold bars as collateral...",
        Category:    "trust-indenture-collateral-agreement",
    },
    Ownership: vc.Ownership{
        Trustor: &vc.Party{
            DID:  "did:dht:EiCFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Xco",
            Name: "Xco",
            Type: "corporation",
            Role: "borrower-trustor",
        },
        Trustee: &vc.Party{
            DID:  "did:dht:EiAFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Abx",
            Name: "Abaxx Technologies Inc.",
            Type: "corporation",
            Role: "lender-trustee",
        },
    },
    // ... other fields
}

// Create credential with specific contexts and types
credential := vc.Create(
    digitalTitle,
    vc.Contexts("https://schemas.abaxx.tech/digital-title/v1"),
    vc.Types("DigitalTitleCredential", "LegalAgreementCredential"),
    vc.Schemas("https://schemas.abaxx.tech/digital-title/schemas/title-record"),
)

// Sign the credential
issuer, _ := didjwk.Create()
vcJWT, err := credential.Sign(issuer)
```

### Verification

```go
// Verify the credential
decoded, err := vc.Verify[*vc.DigitalTitleClaims](vcJWT)
if err != nil {
    // Handle verification error
}

// Access the agreement details
agreement := decoded.VC.CredentialSubject
fmt.Printf("Agreement: %s\n", agreement.Asset.Name)
fmt.Printf("Trustor: %s\n", agreement.Ownership.Trustor.Name)
fmt.Printf("Trustee: %s\n", agreement.Ownership.Trustee.Name)
```

## Approach 2: Reference-Based Claims

This approach creates lightweight credentials that reference the full digital title agreement, suitable for role-based assertions.

### Benefits
- **Compact Size**: Smaller credential size for efficient transmission
- **Role-Specific**: Focus on specific party relationships
- **Evidence-Based**: Link to full agreement via evidence field

### Usage

```go
// Define reference claims struct
type DigitalTitleReferenceClaims struct {
    ID                  string    `json:"id"`
    DigitalTitleID      string    `json:"digitalTitleId"`
    DigitalTitleURI     string    `json:"digitalTitleUri,omitempty"`
    AgreementType       string    `json:"agreementType"`
    PartyRole           string    `json:"partyRole"`
    LegalStatus         string    `json:"legalStatus"`
    EffectiveDate       time.Time `json:"effectiveDate"`
    CollateralValue     *vc.Price `json:"collateralValue,omitempty"`
    AssertionTimestamp  time.Time `json:"assertionTimestamp"`
}

func (d DigitalTitleReferenceClaims) GetID() string { return d.ID }
func (d *DigitalTitleReferenceClaims) SetID(id string) { d.ID = id }

// Create reference claims
titleRef := DigitalTitleReferenceClaims{
    ID:                 "did:dht:EiCFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Xco",
    DigitalTitleID:     "770f9510-e29b-41d4-a716-446655440002",
    DigitalTitleURI:    "https://titles.abaxx.tech/legal-agreements/770f9510-e29b-41d4-a716-446655440002",
    AgreementType:      "trust-indenture-collateral-agreement",
    PartyRole:          "trustor",
    LegalStatus:        "agreement-executed",
    EffectiveDate:      time.Now(),
    CollateralValue: &vc.Price{
        Amount:   2500000.00,
        Currency: "USD",
    },
    AssertionTimestamp: time.Now(),
}

// Create evidence pointing to full agreement
evidence := vc.Evidence{
    ID:   "digital-title-evidence-001",
    Type: "DocumentEvidence",
    AdditionalFields: map[string]interface{}{
        "documentType": "legal-agreement",
        "documentURI":  "https://titles.abaxx.tech/legal-agreements/770f9510-e29b-41d4-a716-446655440002",
        "hashValue":    "sha256:abc123def456...",
        "storageType":  "immutable-ledger",
    },
}

credential := vc.Create(
    &titleRef,
    vc.Types("DigitalTitleParticipantCredential"),
    vc.Contexts("https://schemas.abaxx.tech/digital-title/v1"),
    vc.Evidences(evidence),
)
```

## Approach 3: Generic Claims

This approach uses the flexible `vc.Claims` map for varying structures or rapid prototyping.

### Benefits
- **Flexibility**: Easy to modify structure without code changes
- **Dynamic**: Handle varying legal agreement formats
- **Prototyping**: Quick iteration during development

### Usage

```go
// Create claims using flexible map
claims := vc.Claims{
    "id":            "did:dht:EiCFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Xco",
    "titleId":       "770f9510-e29b-41d4-a716-446655440002",
    "titleType":     "legal-agreement",
    "agreementName": "Gold Collateral Trust Indenture and Margin Lending Agreement",
    "trustorRole":   "borrower-trustor",
    "trusteeRole":   "lender-trustee", 
    "collateralType": "gold-bars",
    "jurisdiction":  "Barbados",
    "executionDate": time.Now().Format(time.RFC3339),
    "keyProvisions": []string{
        "Trust establishment and security interest",
        "Blockchain technology implementation", 
        "Margin lending facility",
        "Default and remedies",
    },
    "digitalRights": map[string]interface{}{
        "tokenization":   "digital-asset-rights-representation",
        "smartContracts": "automated-enforcement-mechanisms",
        "oracles":        "gold-price-feeds-and-vault-verification",
    },
}

credential := vc.Create(
    claims,
    vc.Types("DigitalTitleCredential", "LegalAgreementCredential"),
    vc.Contexts("https://schemas.abaxx.tech/digital-title/v1"),
)
```

## Command Line Usage

The protocol includes a CLI command for creating digital title credentials:

### Full Digital Title Credential

```bash
# Create a full digital title credential
./abaxx-id vc create-digital-title \
    "did:dht:EiCFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Xco" \
    '{"titleId":"770f9510-e29b-41d4-a716-446655440002","titleType":"legal-agreement",...}' \
    --sign '{"did":"did:dht:...","keyManager":{...}}'
```

### Reference Credential

```bash
# Create a reference-based credential
./abaxx-id vc create-digital-title \
    "did:dht:EiCFstcguyQzE90I1uqLRga_i1YgN7urCuf2nESRlV-Xco" \
    '{"titleId":"770f9510-e29b-41d4-a716-446655440002",...}' \
    --use-reference \
    --sign '{"did":"did:dht:...","keyManager":{...}}'
```

## Schema and Context

Digital title credentials use specific JSON-LD contexts and schemas:

### Context
- `https://www.w3.org/2018/credentials/v1` (required base context)
- `https://schemas.abaxx.tech/digital-title/v1` (digital title context)

### Types
- `VerifiableCredential` (required base type)
- `DigitalTitleCredential` (digital title specific)
- `LegalAgreementCredential` (legal agreement specific)
- `DigitalTitleParticipantCredential` (for reference-based credentials)

### Schema
- `https://schemas.abaxx.tech/digital-title/schemas/title-record`

## Security Considerations

### Trust and Verification

1. **Issuer Authority**: Ensure the credential issuer has authority to attest to the legal agreement
2. **Subject Verification**: Verify the credential subject is a legitimate party to the agreement
3. **Temporal Validity**: Check issuance and expiration dates against agreement terms
4. **Evidence Integrity**: Validate evidence hashes and external document references

### Privacy

1. **Selective Disclosure**: Use reference-based approach to limit exposed information
2. **Zero-Knowledge Proofs**: Consider ZKP implementations for sensitive agreement terms
3. **Access Control**: Implement proper access controls for full agreement details

### Legal Considerations

1. **Jurisdiction Compliance**: Ensure credentials comply with applicable legal frameworks
2. **Digital Signatures**: Legal agreements may require specific signature standards
3. **Audit Trail**: Maintain immutable records of credential issuance and usage

## Integration with DWN Protocol

Digital title credentials can be stored and accessed via the Decentralized Web Node (DWN) protocol:

### Storage

```go
// Store credential in DWN
record := dwn.Record{
    Data:     vcJWT,
    Protocol: "https://schemas.abaxx.tech/digital-title/protocol",
    Schema:   "https://schemas.abaxx.tech/digital-title/schemas/credential",
}
```

### Query

```go
// Query for digital title credentials
query := dwn.Query{
    Filter: dwn.Filter{
        Protocol: "https://schemas.abaxx.tech/digital-title/protocol",
        Schema:   "https://schemas.abaxx.tech/digital-title/schemas/credential",
    },
}
```

## Best Practices

1. **Choose the Right Approach**: 
   - Use full credentials for complete legal documentation
   - Use references for role-based assertions
   - Use generic claims for prototyping

2. **Include Evidence**: Always provide evidence linking to authoritative sources

3. **Set Appropriate Expiration**: Align credential expiration with agreement terms

4. **Use Proper Contexts**: Include all necessary JSON-LD contexts for interoperability

5. **Validate Structure**: Ensure compliance with digital title schema

6. **Document Relationships**: Clearly indicate relationships between parties and agreements

## Example Scenarios

### Collateral Lending
A trustor pledges gold bars as collateral for a margin lending facility. The digital title legal agreement is issued as a VC to:
- Establish the trust relationship
- Document collateral details
- Enable automated margin calls via smart contracts

### Asset Transfer
When ownership of a digital title changes hands:
1. Current owner creates a transfer credential
2. New owner receives ownership credential
3. Registry updates the authoritative digital title record

### Regulatory Compliance
Financial institutions can use digital title credentials to:
- Prove compliance with collateral requirements
- Demonstrate legal ownership of assets
- Provide audit trails for regulatory reporting

## Conclusion

Digital title legal agreements can be effectively represented as verifiable credential claims using the Abaxx ID protocol. Choose the appropriate approach based on your use case requirements for completeness, efficiency, and security. 