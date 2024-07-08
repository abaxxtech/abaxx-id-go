package pexv2_test

import (
	"github.com/abaxxtech/abaxx-id-go/internal/pexv2"
)

type PresentationInput struct {
	PresentationDefinition pexv2.PresentationDefinition `json:"presentationDefinition"`
	CredentialJwts         []string                     `json:"credentialJwts"`
}

type PresentationOutput struct {
	SelectedCredentials []string `json:"selectedCredentials"`
}
