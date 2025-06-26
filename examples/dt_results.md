🚀 Digital Title Legal Agreement as Verifiable Credential - Full Example
=======================================================================

📋 Running the complete digital title example...

=== Step 1: Creating DIDs ===
Abaxx DID: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ
Xco DID: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ

=== Step 2: Digital Title Legal Agreement Data ===
Digital Title ID: 770f9510-e29b-41d4-a716-446655440002
Agreement: Gold Collateral Trust Indenture and Margin Lending Agreement

=== Step 3: Creating Verifiable Credential ===
Credential ID: urn:vc:uuid:57230ec5-1fb7-457b-866c-7df7d2425bbf
Issuer: 
Subject: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ

=== Step 4: Signing Credential ===
Signed VC-JWT created successfully!
VC-JWT length: 7553 characters
VC-JWT preview: eyJhbGciOiJFZERTQSIsImtpZCI6ImRpZDpqd2s6ZXlKcmRIa2lPaUpQUzFBaUxDSmpjbllpT2lKRlpESTFOVEU1SWl3aWVDSTZJ...

=== Step 5: Verifying Credential ===
✅ Credential verification successful!
Verified Agreement: Gold Collateral Trust Indenture and Margin Lending Agreement
Verified Trustor: Xco (did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ)
Verified Trustee: Abaxx Technologies Inc. (did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ)
Collateral Type: gold-bars
Jurisdiction: Barbados

=== Step 6: Creating Reference-Based Credential ===
Reference VC-JWT created successfully!
Reference VC-JWT length: 2462 characters
✅ Reference credential verification successful!
Reference Title ID: 770f9510-e29b-41d4-a716-446655440002
Reference Role: trustor

=== Step 7: JSON Representations ===
Full Credential JSON (first 500 chars):
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://schemas.abaxx.tech/digital-title/v1"
  ],
  "type": [
    "VerifiableCredential",
    "DigitalTitleCredential",
    "LegalAgreementCredential"
  ],
  "issuer": "",
  "credentialSubject": {
    "id": "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ",
    "titleId": "770f9510-e29b-41d4-a716-446655440002",
    "titleType": "legal-agreement",
    "asset": {
...

Reference Credential JSON:
{
  "@context": [
    "https://www.w3.org/2018/credentials/v1",
    "https://schemas.abaxx.tech/digital-title/v1"
  ],
  "type": [
    "VerifiableCredential",
    "DigitalTitleParticipantCredential"
  ],
  "issuer": "",
  "credentialSubject": {
    "agreementType": "trust-indenture-collateral-agreement",
    "assertionTimestamp": "2025-06-26T14:33:36-04:00",
    "collateralValue": {
      "amount": 2500000,
      "currency": "USD"
    },
    "digitalTitleId": "770f9510-e29b-41d4-a716-446655440002",
    "digitalTitleUri": "https://titles.abaxx.tech/legal-agreements/770f9510-e29b-41d4-a716-446655440002",
    "effectiveDate": "2024-01-15T10:00:00Z",
    "id": "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ",
    "legalStatus": "agreement-executed",
    "partyRole": "trustor"
  },
  "id": "urn:vc:uuid:31ae3779-7422-4cc5-b372-31fa5fc1d951",
  "issuanceDate": "2025-06-26T18:33:36Z",
  "evidence": [
    {
      "id": "digital-title-evidence-001",
      "type": "DocumentEvidence",
      "AdditionalFields": {
        "documentType": "digital-title-legal-agreement",
        "documentURI": "https://titles.abaxx.tech/legal-agreements/770f9510-e29b-41d4-a716-446655440002",
        "hashValue": "sha256:abc123def456...",
        "storageType": "immutable-ledger",
        "titleId": "770f9510-e29b-41d4-a716-446655440002"
      }
    }
  ]
}

=== Step 8: CLI Command Examples ===
CLI command to create full digital title credential:
./abaxx-id vc create-digital-title \
    "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ" \
    '{"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ","titleId":"770f9510-e29b-41d4-a716-446655440002","titleType":"legal-agreement","asset":{"identifier":"TRUST-INDENTURE-XCO-ABAXX-2024-001","name":"Gold Collateral Trust Indenture and Margin Lending Agreement","description":"Trust arrangement where Xco pledges gold bars as collateral to Abaxx Technologies Inc. for a margin lending facility using blockchain technology","category":"trust-indenture-collateral-agreement","specifications":{"agreementType":"hybrid-security-trust-arrangement","collateralType":"gold-bars","documentLength":"150+ pages","facilityType":"margin-lending-line-of-credit","governingLaw":"Barbados","keyProvisions":["Trust establishment and security interest","Blockchain technology implementation","Margin lending facility","Default and remedies","Digital asset rights management"],"technologyImplementation":"blockchain-smart-contracts"},"location":{"jurisdiction":"Barbados","governingCourt":"Barbados High Court","applicableLaw":"Barbados Commercial Law"}},"ownership":{"trustor":{"did":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ","name":"Xco","type":"corporation","role":"borrower-trustor","contactInfo":{"address":"Xco Corporate Headquarters","email":"legal@xco.com"}},"trustee":{"did":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","name":"Abaxx Technologies Inc.","type":"corporation","role":"lender-trustee","contactInfo":{"address":"Abaxx Technologies Corporate Office","email":"legal@abaxx.com"}},"executionDate":"2024-01-15T10:00:00Z","effectiveDate":"2024-01-15T10:00:00Z"},"legal":{"issuingAuthority":{"name":"Barbados Corporate Registry","jurisdiction":"Barbados","registrationNumber":"TRUST-2024-XCO-ABAXX-001"},"collateralDetails":{"assetType":"gold-bars","digitalRightsManagement":"blockchain-tokens","insuranceRequired":true,"storageLocation":"qualified-vault-facility","valuationMethod":"market-price-oracles"},"restrictions":[{"type":"security-interest","description":"Gold bars pledged as collateral, subject to UCC filings","perfectionMethod":"UCC-1-filing-and-blockchain-registration"},{"type":"margin-requirements","description":"Loan-to-value ratio must be maintained per agreement terms","maintenanceThreshold":"specified-in-line-of-credit-agreement"}],"transferability":{"isTransferable":false,"requiresApproval":true,"approvalAuthority":"both-parties-mutual-consent","restrictions":["Assignment requires written consent of both parties","Subject to regulatory approvals","Must maintain collateral perfection"]}},"valuation":{"collateralValue":{"amount":0,"currency":"USD","lastValuation":"2024-05-15T14:20:00Z","valuationMethod":"real-time-market-pricing","valuationFrequency":"continuous-monitoring","baseAsset":"gold-bars"},"creditFacility":{"currency":"USD","interestRate":"as-per-line-of-credit-agreement","marginCallThreshold":"specified-ltv-ratio","maxCreditLine":"defined-in-schedule-a"}},"technology":{"blockchainImplementation":{"governance":"upgrade-and-dispute-resolution-procedures","keyManagement":"secure-multi-signature-controls","oracles":"gold-price-feeds-and-vault-verification","smartContracts":"automated-enforcement-mechanisms","tokenization":"digital-asset-rights-representation"},"securityProtocols":{"accessControls":"role-based-permissions","auditTrail":"immutable-transaction-records","failsafeMechanisms":"traditional-legal-fallback-procedures"}},"metadata":{"created":"2024-01-15T10:30:00Z","lastUpdated":"2024-05-15T14:20:00Z","version":"1.0","relatedDocuments":[{"documentType":"line-of-credit-agreement","documentId":"schedule-a-loc-agreement","description":"Schedule A - Line of Credit Agreement with commercial terms","relationship":"integrated-facility-terms"},{"documentType":"ucc-filing","documentId":"ucc-1-filing-2024-001","description":"UCC-1 financing statement for security interest perfection","relationship":"legal-perfection"}],"alsoKnownAs":["xco-abaxx-gold-collateral-trust-2024","margin-lending-blockchain-agreement-001","digital-gold-rights-trust-indenture"],"tags":["trust-indenture","gold-collateral","margin-lending","blockchain-implementation","security-agreement","digital-asset-rights","barbados-law","smart-contracts","collateral-management"],"confidentiality":{"level":"confidential","restrictions":"Not to be disseminated outside authorized group"}}}' \
    --sign '{"uri":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","privateKeys":[{"kty":"OKP","crv":"Ed25519","d":"F3PunPwgxLDd3RkLcuQY3P3rSoiVBDzHITO_SjbUTc-qHDkjwYrgfhR6D7zahzJq3Ow9dK5q1w7JGOQ9XLl5KQ","x":"qhw5I8GK4H4Ueg-82ocyatzsPXSuatcOyRjkPVy5eSk"}],"document":{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","verificationMethod":[{"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0","type":"JsonWebKey","controller":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","publicKeyJwk":{"kty":"OKP","crv":"Ed25519","x":"qhw5I8GK4H4Ueg-82ocyatzsPXSuatcOyRjkPVy5eSk"}}],"assertionMethod":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"],"authentication":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"],"capabilityDelegation":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"],"capabilityInvocation":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"]},"metadata":null}'

CLI command to create reference credential:
./abaxx-id vc create-digital-title \
    "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ" \
    '{"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ","titleId":"770f9510-e29b-41d4-a716-446655440002","titleType":"legal-agreement","asset":{"identifier":"TRUST-INDENTURE-XCO-ABAXX-2024-001","name":"Gold Collateral Trust Indenture and Margin Lending Agreement","description":"Trust arrangement where Xco pledges gold bars as collateral to Abaxx Technologies Inc. for a margin lending facility using blockchain technology","category":"trust-indenture-collateral-agreement","specifications":{"agreementType":"hybrid-security-trust-arrangement","collateralType":"gold-bars","documentLength":"150+ pages","facilityType":"margin-lending-line-of-credit","governingLaw":"Barbados","keyProvisions":["Trust establishment and security interest","Blockchain technology implementation","Margin lending facility","Default and remedies","Digital asset rights management"],"technologyImplementation":"blockchain-smart-contracts"},"location":{"jurisdiction":"Barbados","governingCourt":"Barbados High Court","applicableLaw":"Barbados Commercial Law"}},"ownership":{"trustor":{"did":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ","name":"Xco","type":"corporation","role":"borrower-trustor","contactInfo":{"address":"Xco Corporate Headquarters","email":"legal@xco.com"}},"trustee":{"did":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","name":"Abaxx Technologies Inc.","type":"corporation","role":"lender-trustee","contactInfo":{"address":"Abaxx Technologies Corporate Office","email":"legal@abaxx.com"}},"executionDate":"2024-01-15T10:00:00Z","effectiveDate":"2024-01-15T10:00:00Z"},"legal":{"issuingAuthority":{"name":"Barbados Corporate Registry","jurisdiction":"Barbados","registrationNumber":"TRUST-2024-XCO-ABAXX-001"},"collateralDetails":{"assetType":"gold-bars","digitalRightsManagement":"blockchain-tokens","insuranceRequired":true,"storageLocation":"qualified-vault-facility","valuationMethod":"market-price-oracles"},"restrictions":[{"type":"security-interest","description":"Gold bars pledged as collateral, subject to UCC filings","perfectionMethod":"UCC-1-filing-and-blockchain-registration"},{"type":"margin-requirements","description":"Loan-to-value ratio must be maintained per agreement terms","maintenanceThreshold":"specified-in-line-of-credit-agreement"}],"transferability":{"isTransferable":false,"requiresApproval":true,"approvalAuthority":"both-parties-mutual-consent","restrictions":["Assignment requires written consent of both parties","Subject to regulatory approvals","Must maintain collateral perfection"]}},"valuation":{"collateralValue":{"amount":0,"currency":"USD","lastValuation":"2024-05-15T14:20:00Z","valuationMethod":"real-time-market-pricing","valuationFrequency":"continuous-monitoring","baseAsset":"gold-bars"},"creditFacility":{"currency":"USD","interestRate":"as-per-line-of-credit-agreement","marginCallThreshold":"specified-ltv-ratio","maxCreditLine":"defined-in-schedule-a"}},"technology":{"blockchainImplementation":{"governance":"upgrade-and-dispute-resolution-procedures","keyManagement":"secure-multi-signature-controls","oracles":"gold-price-feeds-and-vault-verification","smartContracts":"automated-enforcement-mechanisms","tokenization":"digital-asset-rights-representation"},"securityProtocols":{"accessControls":"role-based-permissions","auditTrail":"immutable-transaction-records","failsafeMechanisms":"traditional-legal-fallback-procedures"}},"metadata":{"created":"2024-01-15T10:30:00Z","lastUpdated":"2024-05-15T14:20:00Z","version":"1.0","relatedDocuments":[{"documentType":"line-of-credit-agreement","documentId":"schedule-a-loc-agreement","description":"Schedule A - Line of Credit Agreement with commercial terms","relationship":"integrated-facility-terms"},{"documentType":"ucc-filing","documentId":"ucc-1-filing-2024-001","description":"UCC-1 financing statement for security interest perfection","relationship":"legal-perfection"}],"alsoKnownAs":["xco-abaxx-gold-collateral-trust-2024","margin-lending-blockchain-agreement-001","digital-gold-rights-trust-indenture"],"tags":["trust-indenture","gold-collateral","margin-lending","blockchain-implementation","security-agreement","digital-asset-rights","barbados-law","smart-contracts","collateral-management"],"confidentiality":{"level":"confidential","restrictions":"Not to be disseminated outside authorized group"}}}' \
    --use-reference \
    --sign '{"uri":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","privateKeys":[{"kty":"OKP","crv":"Ed25519","d":"F3PunPwgxLDd3RkLcuQY3P3rSoiVBDzHITO_SjbUTc-qHDkjwYrgfhR6D7zahzJq3Ow9dK5q1w7JGOQ9XLl5KQ","x":"qhw5I8GK4H4Ueg-82ocyatzsPXSuatcOyRjkPVy5eSk"}],"document":{"@context":["https://www.w3.org/ns/did/v1"],"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","verificationMethod":[{"id":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0","type":"JsonWebKey","controller":"did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ","publicKeyJwk":{"kty":"OKP","crv":"Ed25519","x":"qhw5I8GK4H4Ueg-82ocyatzsPXSuatcOyRjkPVy5eSk"}}],"assertionMethod":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"],"authentication":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"],"capabilityDelegation":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"],"capabilityInvocation":["did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0"]},"metadata":null}'

=== Example Complete! ===
✅ Successfully created and verified digital title credentials
✅ Both full and reference-based approaches demonstrated
✅ CLI commands provided for practical usage

=== Step 9: Cryptographic Proofs ===
🔐 Cryptographic Proof Details:
   Issuer DID: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ
   Subject DID: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ
   JWT ID: urn:vc:uuid:57230ec5-1fb7-457b-866c-7df7d2425bbf
   Signature Algorithm: EdDSA
   Key ID: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ#0
   Signature Length: 64 bytes
   Signature (first 32 bytes): b99433040596b74d1eb61f6927110053ee438c2717541dd385ce703155c837bc...

🔍 Signature Verification:
   ✅ Signature verification successful
   ✅ Credential integrity confirmed
   ✅ Issuer authenticity verified

📋 Creating Verifiable Presentation (VP):
VP Created: urn:vp:uuid:1750962816
VP Holder: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6IkgxZG1HZFFzQmhTdVB5SXFsU2VOVWwzYUxISGw0ZnY5Y0hybUFvUE44TjQifQ
VP contains VC with title: Gold Collateral Trust Indenture and Margin Lending Agreement
VP size: 8320 bytes

🔗 Cryptographic Hash Proof:
   VC-JWT Length: 7553 characters
   Content Hash: simple_a3dd1
   Integrity Status: ✅ Tamper-evident

📊 Cryptographic Proof Summary:
   ✅ Digital Signature: Valid (Ed25519)
   ✅ Issuer Verification: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieC...
   ✅ Subject Authentication: did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieC...
   ✅ Temporal Validity: Current
   ✅ Content Integrity: Verified
   ✅ Schema Compliance: Valid
   ✅ Evidence Chain: Linked

💾 Exportable Proof Package:
Proof data size: 16525 bytes
Proof package preview (first 200 chars):
{
  "agreementHash": "simple_a3dd1",
  "issuerDID": "did:jwk:eyJrdHkiOiJPS1AiLCJjcnYiOiJFZDI1NTE5IiwieCI6InFodzVJOEdLNEg0VWVnLTgyb2N5YXR6c1BYU3VhdGNPeVJqa1BWeTVlU2sifQ",
  "proofType": "Ed25519Signatu...

🎯 Cryptographic Proof Complete!
This proof package can be used to:
   • Verify the digital title legal agreement authenticity
   • Confirm the parties' cryptographic signatures
   • Validate the agreement's integrity and tamper-evidence
   • Enable third-party verification without exposing private keys
   • Support regulatory compliance and audit requirements

🎉 Digital Title Example Complete!

What this example demonstrated:
✅ Created DIDs for trustor (Xco) and trustee (Abaxx Technologies)
✅ Built complete digital title legal agreement structure
✅ Created full verifiable credential with entire agreement
✅ Signed credential to generate VC-JWT
✅ Verified credential and extracted data
✅ Created reference-based credential for role assertions
✅ Generated CLI commands for practical usage
✅ 🔐 Generated cryptographic proofs and verification details
✅ 📋 Created verifiable presentation (VP) for enhanced proof
✅ 🔗 Computed integrity hashes for tamper-evidence
✅ 💾 Exported complete proof package for third-party verification

🔐 CRYPTOGRAPHIC PROOF FEATURES:
   • Ed25519 digital signature verification
   • JWT component analysis (header, payload, signature)
   • Verifiable presentation generation
   • Content integrity hash computation
   • Exportable proof package with verification keys
   • Third-party verification support
