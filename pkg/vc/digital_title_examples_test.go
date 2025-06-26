package vc_test

import (
	"fmt"
	"time"

	"github.com/abaxxtech/abaxx-id-go/pkg/dids/didjwk"
	"github.com/abaxxtech/abaxx-id-go/pkg/vc"
)

// Example_digitalTitleAsFullCredential demonstrates using the entire digital title legal agreement
// as the credential subject
func Example_digitalTitleAsFullCredential() {
	// Create issuer and subject DIDs
	issuer, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	trustorDID, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	// Create a digital title legal agreement claim
	now := time.Now()
	digitalTitle := &vc.DigitalTitleClaims{
		ID:        trustorDID.URI,
		TitleID:   "770f9510-e29b-41d4-a716-446655440002",
		TitleType: "legal-agreement",
		Asset: vc.Asset{
			Identifier:  "TRUST-INDENTURE-XCO-ABAXX-2024-001",
			Name:        "Gold Collateral Trust Indenture and Margin Lending Agreement",
			Description: "Trust arrangement where Xco pledges gold bars as collateral to Abaxx Technologies Inc.",
			Category:    "trust-indenture-collateral-agreement",
			Specifications: map[string]interface{}{
				"agreementType":            "hybrid-security-trust-arrangement",
				"collateralType":           "gold-bars",
				"facilityType":             "margin-lending-line-of-credit",
				"technologyImplementation": "blockchain-smart-contracts",
				"governingLaw":             "Barbados",
			},
			Location: &vc.Location{
				Jurisdiction:   "Barbados",
				GoverningCourt: "Barbados High Court",
				ApplicableLaw:  "Barbados Commercial Law",
			},
		},
		Ownership: vc.Ownership{
			Trustor: &vc.Party{
				DID:  trustorDID.URI,
				Name: "Xco",
				Type: "corporation",
				Role: "borrower-trustor",
			},
			Trustee: &vc.Party{
				DID:  issuer.URI,
				Name: "Abaxx Technologies Inc.",
				Type: "corporation",
				Role: "lender-trustee",
			},
			ExecutionDate: &now,
			EffectiveDate: &now,
		},
		Legal: vc.Legal{
			IssuingAuthority: vc.IssuingAuthority{
				Name:               "Barbados Corporate Registry",
				Jurisdiction:       "Barbados",
				RegistrationNumber: "TRUST-2024-XCO-ABAXX-001",
			},
			CollateralDetails: map[string]interface{}{
				"assetType":               "gold-bars",
				"storageLocation":         "qualified-vault-facility",
				"digitalRightsManagement": "blockchain-tokens",
				"valuationMethod":         "market-price-oracles",
				"insuranceRequired":       true,
			},
			Restrictions: []vc.Restriction{
				{
					Type:             "security-interest",
					Description:      "Gold bars pledged as collateral, subject to UCC filings",
					PerfectionMethod: "UCC-1-filing-and-blockchain-registration",
				},
			},
			Transferability: vc.Transferability{
				IsTransferable:    false,
				RequiresApproval:  true,
				ApprovalAuthority: "both-parties-mutual-consent",
			},
		},
		Technology: &vc.Technology{
			BlockchainImplementation: map[string]interface{}{
				"tokenization":   "digital-asset-rights-representation",
				"smartContracts": "automated-enforcement-mechanisms",
				"oracles":        "gold-price-feeds-and-vault-verification",
				"keyManagement":  "secure-multi-signature-controls",
			},
		},
		Metadata: vc.Metadata{
			Created:     now,
			LastUpdated: now,
			Version:     "1.0",
			Tags:        []string{"trust-indenture", "gold-collateral", "margin-lending", "blockchain-implementation"},
		},
	}

	// Create the verifiable credential with digital title context and type
	credential := vc.Create(
		digitalTitle,
		vc.Contexts("https://schemas.abaxx.tech/digital-title/v1"),
		vc.Types("DigitalTitleCredential", "LegalAgreementCredential"),
		vc.Schemas("https://schemas.abaxx.tech/digital-title/schemas/title-record"),
	)

	// Sign the credential
	vcJWT, err := credential.Sign(issuer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Digital Title VC-JWT created successfully")

	// Verify the credential
	decoded, err := vc.Verify[*vc.DigitalTitleClaims](vcJWT)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Agreement: %s\n", decoded.VC.CredentialSubject.Asset.Name)
	fmt.Printf("Trustor: %s\n", decoded.VC.CredentialSubject.Ownership.Trustor.Name)
	fmt.Printf("Trustee: %s\n", decoded.VC.CredentialSubject.Ownership.Trustee.Name)

	// Output: Digital Title VC-JWT created successfully
	// Agreement: Gold Collateral Trust Indenture and Margin Lending Agreement
	// Trustor: Xco
	// Trustee: Abaxx Technologies Inc.
}

// DigitalTitleReferenceClaims represents a credential that references a digital title
type DigitalTitleReferenceClaims struct {
	ID                 string     `json:"id"`
	DigitalTitleID     string     `json:"digitalTitleId"`
	DigitalTitleURI    string     `json:"digitalTitleUri,omitempty"`
	AgreementType      string     `json:"agreementType"`
	PartyRole          string     `json:"partyRole"`
	LegalStatus        string     `json:"legalStatus"`
	EffectiveDate      time.Time  `json:"effectiveDate"`
	ExpirationDate     *time.Time `json:"expirationDate,omitempty"`
	CollateralValue    *vc.Price  `json:"collateralValue,omitempty"`
	AssertionTimestamp time.Time  `json:"assertionTimestamp"`
}

func (d DigitalTitleReferenceClaims) GetID() string {
	return d.ID
}

func (d *DigitalTitleReferenceClaims) SetID(id string) {
	d.ID = id
}

// Example_digitalTitleAsReference demonstrates using a reference to the digital title
// with specific claims about the party's relationship to the agreement
func Example_digitalTitleAsReference() {
	// Create issuer and subject DIDs
	issuer, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	trustorDID, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	now := time.Now()

	// Create reference claims to the digital title legal agreement
	titleRef := DigitalTitleReferenceClaims{
		ID:              trustorDID.URI,
		DigitalTitleID:  "770f9510-e29b-41d4-a716-446655440002",
		DigitalTitleURI: "https://titles.abaxx.tech/legal-agreements/770f9510-e29b-41d4-a716-446655440002",
		AgreementType:   "trust-indenture-collateral-agreement",
		PartyRole:       "trustor",
		LegalStatus:     "agreement-executed",
		EffectiveDate:   now,
		CollateralValue: &vc.Price{
			Amount:   2500000.00,
			Currency: "USD",
		},
		AssertionTimestamp: now,
	}

	// Create credential with additional evidence pointing to the full agreement
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

	// Sign the credential
	vcJWT, err := credential.Sign(issuer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Digital Title Reference VC-JWT created successfully")

	// Verify the credential
	decoded, err := vc.Verify[*DigitalTitleReferenceClaims](vcJWT)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Title ID: %s\n", decoded.VC.CredentialSubject.DigitalTitleID)
	fmt.Printf("Role: %s\n", decoded.VC.CredentialSubject.PartyRole)
	fmt.Printf("Status: %s\n", decoded.VC.CredentialSubject.LegalStatus)

	// Output: Digital Title Reference VC-JWT created successfully
	// Title ID: 770f9510-e29b-41d4-a716-446655440002
	// Role: trustor
	// Status: agreement-executed
}

// Example_digitalTitleWithGenericClaims demonstrates using the flexible Claims approach
// for scenarios where the structure varies
func Example_digitalTitleWithGenericClaims() {
	// Create issuer and subject DIDs
	issuer, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	trustorDID, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	// Create claims using the flexible Claims map
	claims := vc.Claims{
		"id":             trustorDID.URI,
		"titleId":        "770f9510-e29b-41d4-a716-446655440002",
		"titleType":      "legal-agreement",
		"agreementName":  "Gold Collateral Trust Indenture and Margin Lending Agreement",
		"trustorRole":    "borrower-trustor",
		"trusteeRole":    "lender-trustee",
		"collateralType": "gold-bars",
		"jurisdiction":   "Barbados",
		"executionDate":  time.Now().Format(time.RFC3339),
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

	// Sign the credential
	vcJWT, err := credential.Sign(issuer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generic Digital Title VC-JWT created successfully")

	// Verify using generic Claims type
	decoded, err := vc.Verify[vc.Claims](vcJWT)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Agreement: %s\n", decoded.VC.CredentialSubject["agreementName"])
	fmt.Printf("Collateral: %s\n", decoded.VC.CredentialSubject["collateralType"])
	fmt.Printf("Jurisdiction: %s\n", decoded.VC.CredentialSubject["jurisdiction"])

	// Output: Generic Digital Title VC-JWT created successfully
	// Agreement: Gold Collateral Trust Indenture and Margin Lending Agreement
	// Collateral: gold-bars
	// Jurisdiction: Barbados
}
