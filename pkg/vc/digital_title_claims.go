package vc

import (
	"time"
)

// DigitalTitleClaims represents a digital title legal agreement as a verifiable credential subject
type DigitalTitleClaims struct {
	ID         string      `json:"id"` // Required for CredentialSubject interface
	TitleID    string      `json:"titleId"`
	TitleType  string      `json:"titleType"`
	Asset      Asset       `json:"asset"`
	Ownership  Ownership   `json:"ownership"`
	Legal      Legal       `json:"legal"`
	Valuation  *Valuation  `json:"valuation,omitempty"`
	Technology *Technology `json:"technology,omitempty"`
	Metadata   Metadata    `json:"metadata"`
}

// GetID implements CredentialSubject interface
func (d DigitalTitleClaims) GetID() string {
	return d.ID
}

// SetID implements CredentialSubject interface
func (d *DigitalTitleClaims) SetID(id string) {
	d.ID = id
}

// Asset represents the asset information in a digital title
type Asset struct {
	Identifier     string                 `json:"identifier"`
	Name           string                 `json:"name"`
	Description    string                 `json:"description"`
	Category       string                 `json:"category"`
	Specifications map[string]interface{} `json:"specifications,omitempty"`
	Location       *Location              `json:"location,omitempty"`
}

// Location represents the physical or legal location of an asset
type Location struct {
	Address        string       `json:"address,omitempty"`
	City           string       `json:"city,omitempty"`
	State          string       `json:"state,omitempty"`
	Country        string       `json:"country,omitempty"`
	PostalCode     string       `json:"postalCode,omitempty"`
	Coordinates    *Coordinates `json:"coordinates,omitempty"`
	Jurisdiction   string       `json:"jurisdiction,omitempty"`
	GoverningCourt string       `json:"governingCourt,omitempty"`
	ApplicableLaw  string       `json:"applicableLaw,omitempty"`
}

// Coordinates represents geographic coordinates
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Ownership represents ownership information
type Ownership struct {
	Owner             *Party     `json:"owner,omitempty"`
	Trustor           *Party     `json:"trustor,omitempty"`
	Trustee           *Party     `json:"trustee,omitempty"`
	CoOwners          []CoOwner  `json:"coOwners,omitempty"`
	AcquisitionDate   *time.Time `json:"acquisitionDate,omitempty"`
	ExecutionDate     *time.Time `json:"executionDate,omitempty"`
	EffectiveDate     *time.Time `json:"effectiveDate,omitempty"`
	AcquisitionMethod string     `json:"acquisitionMethod,omitempty"`
	PurchasePrice     *Price     `json:"purchasePrice,omitempty"`
}

// Party represents a party in a legal agreement
type Party struct {
	DID         string                 `json:"did"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Role        string                 `json:"role,omitempty"`
	ContactInfo map[string]interface{} `json:"contactInfo,omitempty"`
}

// CoOwner represents a co-owner in joint ownership
type CoOwner struct {
	DID                 string  `json:"did"`
	Name                string  `json:"name"`
	OwnershipPercentage float64 `json:"ownershipPercentage,omitempty"`
}

// Price represents monetary value
type Price struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// Legal represents legal information and restrictions
type Legal struct {
	IssuingAuthority  IssuingAuthority       `json:"issuingAuthority"`
	CollateralDetails map[string]interface{} `json:"collateralDetails,omitempty"`
	Restrictions      []Restriction          `json:"restrictions,omitempty"`
	Transferability   Transferability        `json:"transferability"`
}

// IssuingAuthority represents the authority that issued/recognizes the title
type IssuingAuthority struct {
	Name               string `json:"name"`
	Jurisdiction       string `json:"jurisdiction"`
	DID                string `json:"did,omitempty"`
	RegistrationNumber string `json:"registrationNumber,omitempty"`
}

// Restriction represents legal restrictions on the title
type Restriction struct {
	Type                 string     `json:"type"`
	Description          string     `json:"description"`
	ExpirationDate       *time.Time `json:"expirationDate,omitempty"`
	PerfectionMethod     string     `json:"perfectionMethod,omitempty"`
	MaintenanceThreshold string     `json:"maintenanceThreshold,omitempty"`
}

// Transferability represents rules governing transfer of the title
type Transferability struct {
	IsTransferable    bool     `json:"isTransferable"`
	RequiresApproval  bool     `json:"requiresApproval,omitempty"`
	ApprovalAuthority string   `json:"approvalAuthority,omitempty"`
	Restrictions      []string `json:"restrictions,omitempty"`
}

// Valuation represents current valuation information
type Valuation struct {
	AssessedValue   *ValueInfo             `json:"assessedValue,omitempty"`
	MarketValue     *ValueInfo             `json:"marketValue,omitempty"`
	CollateralValue *ValueInfo             `json:"collateralValue,omitempty"`
	CreditFacility  map[string]interface{} `json:"creditFacility,omitempty"`
}

// ValueInfo represents valuation details
type ValueInfo struct {
	Amount             float64    `json:"amount"`
	Currency           string     `json:"currency"`
	AssessmentDate     *time.Time `json:"assessmentDate,omitempty"`
	ValuationDate      *time.Time `json:"valuationDate,omitempty"`
	LastValuation      *time.Time `json:"lastValuation,omitempty"`
	Assessor           string     `json:"assessor,omitempty"`
	ValuationMethod    string     `json:"valuationMethod,omitempty"`
	ValuationFrequency string     `json:"valuationFrequency,omitempty"`
	BaseAsset          string     `json:"baseAsset,omitempty"`
}

// Technology represents blockchain and technology implementation details
type Technology struct {
	BlockchainImplementation map[string]interface{} `json:"blockchainImplementation,omitempty"`
	SecurityProtocols        map[string]interface{} `json:"securityProtocols,omitempty"`
}

// Metadata represents additional metadata and references
type Metadata struct {
	Created              time.Time         `json:"created"`
	LastUpdated          time.Time         `json:"lastUpdated"`
	Version              string            `json:"version"`
	PrecedingTitleID     string            `json:"precedingTitleId,omitempty"`
	PrecedingAgreementID string            `json:"precedingAgreementId,omitempty"`
	RelatedDocuments     []RelatedDocument `json:"relatedDocuments,omitempty"`
	Sections             []DocumentSection `json:"sections,omitempty"`
	AlsoKnownAs          []string          `json:"alsoKnownAs,omitempty"`
	Tags                 []string          `json:"tags,omitempty"`
	Confidentiality      *Confidentiality  `json:"confidentiality,omitempty"`
}

// RelatedDocument represents a reference to a supporting document
type RelatedDocument struct {
	DocumentType string `json:"documentType"`
	DocumentID   string `json:"documentId"`
	Description  string `json:"description"`
	URL          string `json:"url,omitempty"`
	Relationship string `json:"relationship,omitempty"`
}

// DocumentSection represents a section in a legal document
type DocumentSection struct {
	Number      int    `json:"number"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Confidentiality represents confidentiality information
type Confidentiality struct {
	Level          string `json:"level"`
	Restrictions   string `json:"restrictions"`
	AccessPassword string `json:"accessPassword,omitempty"`
}
