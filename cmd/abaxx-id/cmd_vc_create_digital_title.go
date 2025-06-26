package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/abaxxtech/abaxx-id-go/pkg/dids/did"
	"github.com/abaxxtech/abaxx-id-go/pkg/vc"
)

type vcCreateDigitalTitleCMD struct {
	CredentialSubjectID string    `arg:"" help:"The Credential Subject's ID (DID)"`
	DigitalTitleJSON    string    `arg:"" help:"Digital title legal agreement as JSON string"`
	Sign                string    `help:"Portable DID used to sign the VC-JWT. Value is a JSON string."`
	UseReference        bool      `help:"Create a reference credential instead of embedding full title" default:"false"`
	Contexts            []string  `help:"Add additional @context's to the default ones."`
	Types               []string  `help:"Add additional type's to the default ones."`
	ID                  string    `help:"Override the default ID of format urn:vc:uuid:<uuid>."`
	IssuanceDate        time.Time `help:"Override the default issuanceDate of time.Now()."`
	ExpirationDate      time.Time `help:"Override the default expirationDate of nil."`
	NoIndent            bool      `help:"Print the VC without indentation." default:"false"`
}

func (c *vcCreateDigitalTitleCMD) Run() error {
	opts := []vc.CreateOption{}

	// Set default contexts for digital title credentials
	defaultContexts := []string{"https://schemas.abaxx.tech/digital-title/v1"}
	if len(c.Contexts) > 0 {
		defaultContexts = append(defaultContexts, c.Contexts...)
	}
	opts = append(opts, vc.Contexts(defaultContexts...))

	// Set default types for digital title credentials
	defaultTypes := []string{"DigitalTitleCredential", "LegalAgreementCredential"}
	if len(c.Types) > 0 {
		defaultTypes = append(defaultTypes, c.Types...)
	}
	opts = append(opts, vc.Types(defaultTypes...))

	// Add schema reference
	opts = append(opts, vc.Schemas("https://schemas.abaxx.tech/digital-title/schemas/title-record"))

	if c.ID != "" {
		opts = append(opts, vc.ID(c.ID))
	}
	if (c.IssuanceDate != time.Time{}) {
		opts = append(opts, vc.IssuanceDate(c.IssuanceDate))
	}
	if (c.ExpirationDate != time.Time{}) {
		opts = append(opts, vc.ExpirationDate(c.ExpirationDate))
	}

	if c.UseReference {
		// Create a reference credential
		var rawTitle map[string]interface{}
		err := json.Unmarshal([]byte(c.DigitalTitleJSON), &rawTitle)
		if err != nil {
			return fmt.Errorf("invalid digital title JSON: %w", err)
		}

		// Extract key information for the reference
		titleRef := map[string]interface{}{
			"id":                 c.CredentialSubjectID,
			"digitalTitleId":     rawTitle["titleId"],
			"agreementType":      rawTitle["titleType"],
			"assertionTimestamp": time.Now().Format(time.RFC3339),
		}

		// Add specific fields based on title type
		if asset, ok := rawTitle["asset"].(map[string]interface{}); ok {
			titleRef["agreementName"] = asset["name"]
			titleRef["assetCategory"] = asset["category"]
		}

		if ownership, ok := rawTitle["ownership"].(map[string]interface{}); ok {
			// Determine the party's role
			if trustor, ok := ownership["trustor"].(map[string]interface{}); ok {
				if trustorDID, ok := trustor["did"].(string); ok && trustorDID == c.CredentialSubjectID {
					titleRef["partyRole"] = "trustor"
				}
			}
			if trustee, ok := ownership["trustee"].(map[string]interface{}); ok {
				if trusteeDID, ok := trustee["did"].(string); ok && trusteeDID == c.CredentialSubjectID {
					titleRef["partyRole"] = "trustee"
				}
			}
		}

		claims := vc.Claims(titleRef)

		// Add evidence pointing to the full digital title
		evidence := vc.Evidence{
			ID:   "digital-title-evidence",
			Type: "DocumentEvidence",
			AdditionalFields: map[string]interface{}{
				"documentType": "digital-title-legal-agreement",
				"titleId":      rawTitle["titleId"],
				"fullDocument": rawTitle,
			},
		}

		// Create credential with evidence
		credential := vc.Create(claims, append(opts, vc.Evidences(evidence))...)

		if c.Sign != "" {
			var portableDID did.PortableDID
			err := json.Unmarshal([]byte(c.Sign), &portableDID)
			if err != nil {
				return fmt.Errorf("invalid portable DID: %w", err)
			}

			bearerDID, err := did.FromPortableDID(portableDID)
			if err != nil {
				return err
			}

			signed, err := credential.Sign(bearerDID)
			if err != nil {
				return err
			}

			fmt.Println(signed)
			return nil
		}

		var jsonVC []byte
		if c.NoIndent {
			jsonVC, err = json.Marshal(credential)
		} else {
			jsonVC, err = json.MarshalIndent(credential, "", "  ")
		}
		if err != nil {
			return err
		}

		fmt.Println(string(jsonVC))
		return nil

	} else {
		// Parse the digital title into strongly typed claims
		var digitalTitle vc.DigitalTitleClaims
		err := json.Unmarshal([]byte(c.DigitalTitleJSON), &digitalTitle)
		if err != nil {
			return fmt.Errorf("invalid digital title JSON: %w", err)
		}

		digitalTitle.ID = c.CredentialSubjectID
		credential := vc.Create(&digitalTitle, opts...)

		if c.Sign != "" {
			var portableDID did.PortableDID
			err := json.Unmarshal([]byte(c.Sign), &portableDID)
			if err != nil {
				return fmt.Errorf("invalid portable DID: %w", err)
			}

			bearerDID, err := did.FromPortableDID(portableDID)
			if err != nil {
				return err
			}

			signed, err := credential.Sign(bearerDID)
			if err != nil {
				return err
			}

			fmt.Println(signed)
			return nil
		}

		var jsonVC []byte
		if c.NoIndent {
			jsonVC, err = json.Marshal(credential)
		} else {
			jsonVC, err = json.MarshalIndent(credential, "", "  ")
		}
		if err != nil {
			return err
		}

		fmt.Println(string(jsonVC))
		return nil
	}
}
