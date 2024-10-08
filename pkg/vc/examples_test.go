package vc_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/abaxxtech/abaxx-id-go/pkg/dids/didjwk"
	"github.com/abaxxtech/abaxx-id-go/pkg/vc"
)

// Demonstrates how to create, sign, and verify a Verifiable Credential using the vc package.
func Example() {
	// create sample issuer and subject DIDs
	issuer, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	subject, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	// creation
	claims := vc.Claims{"id": subject.URI, "name": "John Doe"}
	cred := vc.Create(claims)

	// signing
	vcJWT, err := cred.Sign(issuer)
	if err != nil {
		panic(err)
	}

	// verification
	decoded, err := vc.Verify[vc.Claims](vcJWT)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoded.VC.CredentialSubject["name"])
	// Output: John Doe
}

type KnownCustomerClaims struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c KnownCustomerClaims) GetID() string {
	return c.ID
}

func (c *KnownCustomerClaims) SetID(id string) {
	c.ID = id
}

// Demonstrates how to use a strongly typed credential subject
func Example_stronglyTyped() {
	issuer, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	subject, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	claims := KnownCustomerClaims{ID: subject.URI, Name: "John Doe"}
	cred := vc.Create(&claims)

	vcJWT, err := cred.Sign(issuer)
	if err != nil {
		panic(err)
	}

	decoded, err := vc.Verify[*KnownCustomerClaims](vcJWT)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoded.VC.CredentialSubject.Name)
	// Output: John Doe
}

// Demonstrates how to use a mix of strongly typed and untyped credential subjects with the vc package.
func Example_mixed() {
	issuer, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	subject, err := didjwk.Create()
	if err != nil {
		panic(err)
	}

	claims := KnownCustomerClaims{ID: subject.URI, Name: "John Doe"}
	cred := vc.Create(&claims)

	vcJWT, err := cred.Sign(issuer)
	if err != nil {
		panic(err)
	}

	decoded, err := vc.Verify[vc.Claims](vcJWT)
	if err != nil {
		panic(err)
	}

	fmt.Println(decoded.VC.CredentialSubject["name"])
	// Output: John Doe
}

// Demonstrates how to create a Verifiable Credential
func ExampleCreate() {
	claims := vc.Claims{"name": "John Doe"}
	cred := vc.Create(claims)

	bytes, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

// Demonstrates how to create a Verifiable Credential with options
func ExampleCreate_options() {
	claims := vc.Claims{"id": "1234"}
	issuanceDate := time.Now().UTC().Add(10 * time.Hour)
	expirationDate := issuanceDate.Add(30 * time.Hour)

	cred := vc.Create(
		claims,
		vc.ID("thecustomid"),
		vc.Contexts("https://nocontextisbestcontext.gov"),
		vc.Types("StreetCredential"),
		vc.IssuanceDate(issuanceDate),
		vc.ExpirationDate(expirationDate),
	)

	bytes, err := json.MarshalIndent(cred, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}
