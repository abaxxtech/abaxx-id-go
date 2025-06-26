# Digital Title Legal Agreement Examples

This directory contains comprehensive examples of how to use digital title legal agreements as verifiable credential claims in the Abaxx ID protocol.

## 🎯 What's Included

### `digital_title_full_example.go`
A complete end-to-end example demonstrating:
- **DID Creation**: Generate DIDs for all parties (trustor, trustee)
- **Digital Title Structure**: Build the complete legal agreement using strongly-typed Go structs
- **Credential Creation**: Create W3C-compliant verifiable credentials
- **Signing & Verification**: Sign with DIDs and verify cryptographically
- **Multiple Approaches**: Show both full and reference-based credential patterns
- **CLI Integration**: Generate practical command-line examples

### `run_digital_title_example.sh`
Simple script to execute the full example with proper setup and explanation.

## 🚀 Quick Start

```bash
# Make the script executable (if not already)
chmod +x run_digital_title_example.sh

# Run the complete example
./run_digital_title_example.sh
```

Or run directly with Go:

```bash
go run digital_title_full_example.go
```

## 📋 Example Output

The example demonstrates:

1. **DID Generation**
   ```
   Abaxx DID: did:jwk:eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSIs...
   Xco DID: did:jwk:eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSIs...
   ```

2. **Digital Title Structure**
   ```
   Agreement: Gold Collateral Trust Indenture and Margin Lending Agreement
   Trustor: Xco (borrower-trustor)
   Trustee: Abaxx Technologies Inc. (lender-trustee)
   ```

3. **Credential Creation & Signing**
   ```
   ✅ Credential verification successful!
   ✅ Reference credential verification successful!
   ```

4. **Cryptographic Proof Generation**
   ```
   🔐 Cryptographic Proof Details:
      Issuer DID: did:jwk:eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSIs...
      Signature Algorithm: EdDSA
      ✅ Signature verification successful
      ✅ Credential integrity confirmed
   
   📋 Creating Verifiable Presentation (VP):
      VP Created: urn:vp:uuid:1735090234
      VP contains VC with title: Gold Collateral Trust Indenture...
   
   💾 Exportable Proof Package:
      Proof data size: 8,432 bytes
   ```

5. **CLI Commands**
   ```bash
   ./abaxx-id vc create-digital-title \
       "did:jwk:..." \
       '{"titleId":"770f9510-e29b-41d4-a716-446655440002",...}' \
       --sign '{"did":"...","keyManager":{...}}'
   ```

## 🏗️ Digital Title Structure

The example uses the complete digital title legal agreement structure:

```go
type DigitalTitleClaims struct {
    ID        string     // DID of the credential subject
    TitleID   string     // Unique identifier for the title
    TitleType string     // Type: "legal-agreement"
    Asset     Asset      // Complete asset information
    Ownership Ownership  // Trustor/Trustee details with DIDs
    Legal     Legal      // Legal restrictions and transferability
    Valuation *Valuation // Collateral and credit facility values
    Technology *Technology // Blockchain implementation details
    Metadata  Metadata   // Document sections, tags, confidentiality
}
```

## 🔍 Key Features Demonstrated

### 1. **Full Digital Title Credential**
- Embeds the complete legal agreement as credential subject
- Uses strongly-typed Go structs for data integrity
- Includes all agreement details (parties, collateral, technology)

### 2. **Reference-Based Credential**
- Lightweight credential for role-based assertions
- Links to full agreement via evidence field
- Suitable for efficient transmission and storage

### 3. **W3C Compliance**
- Standard verifiable credential format
- Proper JSON-LD contexts and types
- Schema validation support

### 4. **DID Integration**
- Uses `did:jwk` for party identification
- Cryptographic signing and verification
- Portable DID format for CLI usage

### 5. **Blockchain Integration Ready**
- Technology field for smart contract details
- Oracle integration for real-time valuations
- Audit trail and access control specifications

### 6. **Cryptographic Proof Generation**
- **JWT Component Analysis**: Extract and verify signature components
- **Verifiable Presentations**: Create VPs for enhanced proof delivery
- **Integrity Verification**: Generate tamper-evident content hashes
- **Exportable Proof Packages**: Bundle all verification data for third parties
- **Ed25519 Signatures**: Industry-standard cryptographic verification

## 🎨 Use Cases

### **Collateral Lending**
```go
// Trustor pledges gold bars as collateral
digitalTitle.Legal.CollateralDetails = map[string]interface{}{
    "assetType":       "gold-bars",
    "storageLocation": "qualified-vault-facility",
    "valuationMethod": "market-price-oracles",
}
```

### **Regulatory Compliance**
```go
// Maintain audit trails and legal restrictions
digitalTitle.Legal.Restrictions = []vc.Restriction{
    {
        Type:        "security-interest",
        Description: "Gold bars pledged as collateral, subject to UCC filings",
    },
}
```

### **Smart Contract Integration**
```go
// Enable automated enforcement
digitalTitle.Technology.BlockchainImplementation = map[string]interface{}{
    "smartContracts": "automated-enforcement-mechanisms",
    "oracles":        "gold-price-feeds-and-vault-verification",
}
```

### **Cryptographic Proof for Auditing**
```go
// Generate verifiable proof package
proofData := map[string]interface{}{
    "vcJWT":           vcJWT,
    "verificationKey": issuerDID.Document.VerificationMethod[0].PublicKeyJwk,
    "agreementHash":   contentHash,
    "proofType":       "Ed25519Signature2018",
    "verifiablePresentation": vp,
}
```

## 🔧 CLI Integration

The example generates ready-to-use CLI commands:

### Create Full Credential
```bash
./abaxx-id vc create-digital-title \
    "did:jwk:eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSIs..." \
    '{"titleId":"770f9510-e29b-41d4-a716-446655440002","titleType":"legal-agreement",...}' \
    --sign '{"did":"did:jwk:...","keyManager":{...}}'
```

### Create Reference Credential
```bash
./abaxx-id vc create-digital-title \
    "did:jwk:eyJ0eXAiOiJKV1QiLCJhbGciOiJFZERTQSIs..." \
    '{"titleId":"770f9510-e29b-41d4-a716-446655440002",...}' \
    --use-reference \
    --sign '{"did":"did:jwk:...","keyManager":{...}}'
```

## 🔐 Security Features

- **Cryptographic Signatures**: All credentials signed with DID keys (Ed25519)
- **Temporal Validity**: Proper issuance and expiration dates
- **Evidence Integrity**: Hash-based links to authoritative sources
- **Access Control**: Confidentiality levels and restrictions
- **Audit Trail**: Immutable record of credential lifecycle
- **Proof Generation**: Exportable verification packages for third parties
- **Tamper Detection**: Content integrity hashes for tamper-evidence
- **Verifiable Presentations**: Enhanced proof delivery and selective disclosure

## 🚀 Next Steps

1. **Run the Example**: Start with `./run_digital_title_example.sh`
2. **Modify for Your Use Case**: Adapt the legal agreement structure
3. **Integrate with DWN**: Store and query credentials
4. **Add Business Logic**: Implement margin calls and compliance checks
5. **Deploy Smart Contracts**: Enable automated enforcement

## 📚 Related Documentation

- [`../pkg/vc/README_DIGITAL_TITLE.md`](../pkg/vc/README_DIGITAL_TITLE.md) - Comprehensive guide
- [`../json-schemas/digital-title-record.json`](../json-schemas/digital-title-record.json) - Schema definition
- [`../json-schemas/examples/`](../json-schemas/examples/) - Sample data structures

---

**Happy coding! 🎉** The digital title legal agreement becomes a first-class verifiable claim that enables secure, auditable, and automated legal processes. 