package pexv2_test

import (
	"testing"

	"github.com/abaxxtech/abaxx-id-go/pkg/pexv2"
)

type PresentationInput struct {
	PresentationDefinition pexv2.PresentationDefinition `json:"presentationDefinition"`
	CredentialJwts         []string                     `json:"credentialJwts"`
}

type PresentationOutput struct {
	SelectedCredentials []string `json:"selectedCredentials"`
}

func TestPresentation(t *testing.T) {
	t.Skip("skipping")
}
