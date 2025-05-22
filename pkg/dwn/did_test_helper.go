package dwn

import (
	"errors"
	"slices"
	"strings"
)

// DID Resolver for test purposes
type StaticDidResolver struct {
	method   string
	did      string // the singular did we can resolve
	document interface{}
}

func (o StaticDidResolver) Method() string {
	return o.method
}

func (o StaticDidResolver) Resolve(did string) (DidResolutionResult, error) {
	if did != o.did {
		return DidResolutionResult{},
			errors.New("Static DID only knows" + o.did)
	}

	return DidResolutionResult{DidDocument: o.document}, nil
}

func NewStaticDidResolver(didjson string) *DidResolver {
	return NewDidResolver([]DidMethodResolver{StaticDidResolver{
		method:   "ion",
		did:      "foobar",
		document: "not yet!",
	}},
		nil,
	)
}

// This is the entrypoint for the DID output from our API we use in tests.
type fileDID struct {
	Id        string `json:"_id"`
	Did       string `json:"did"`
	Timestamp int64  `json:"timestamp"`

	Metadata didMetadata `json:"metadata"`
}

func (did fileDID) GetKeyById(keyid string) (interface{}, error) {
	parts := strings.Split(keyid, "#")

	if len(parts) != 2 {
		return nil, errors.New("Keyid doesn't have # in it")
	}
	if parts[0] != "" {
		if parts[0] != did.Did {
			return nil, errors.New("KeyID not found: base DID doesnt match")
		}
	}

	keyidtofind := "#" + parts[1]

	// retrieve the keys:
	didDocument := did.Metadata.Did.Document
	keys := slices.Concat(didDocument.KeyAgreement, didDocument.Authentication,
		didDocument.VerificationMethod)

	for _, key := range keys {
		switch v := key.(type) {
		case map[string]interface{}:
			//fmt.Println("iterating keys", v)
			dockeyid, ok := v["id"]
			if ok && dockeyid == keyidtofind {
				jwk, ok := v["publicKeyJwk"]
				if !ok {
					return nil, errors.New("KeyID doesnt have public key in document!")
				}
				return jwk, nil
			}
		}
	}
	return nil, errors.New("KeyID not found")
}

type didMetadata struct {
	// Scalar fields
	Email             string `json:"email"`
	DeviceId          string `json:"deviceId"`
	Identity          string `json:"identity"`
	Username          string `json:"username"`
	CreatedAt         int64  `json:"createdAt"`
	UpdatedAt         int64  `json:"updatedAt"`
	AccountIcon       string `json:"accountIcon"`
	AccountLabel      string `json:"accountLabel"`
	IsFirstTimeLogin  bool   `json:"isFirstTimeLogin"`
	PrimaryAccountId  string `json:"primaryAccountId"`
	VerifierPushToken string `json:"verifierPushToken"`

	// Contains keys, did document, more.
	Did struct {
		Did         string      `json:"did"`
		KeySet      keySet      `json:"keySet"`
		Document    didDocument `json:"document"`
		CanonicalId string      `json:"canonicalId"`
	} `json:"did"`
}

type keySet struct {
	UpdateKey struct {
		PublicKeyJwk  interface{} `json:"publicKeyJwk"`
		PrivateKeyJwk interface{} `json:"privateKeyJwk"`
	} `json:"updateKey"`
	RecoveryKey struct {
		PublicKeyJwk  interface{} `json:"publicKeyJwk"`
		PrivateKeyJwk interface{} `json:"privateKeyJwk"`
	} `json:"recoveryKey"`
	VerificationMethodKeys []struct {
		PublicKeyJwk  interface{} `json:"publicKeyJwk"`
		PrivateKeyJwk interface{} `json:"privateKeyJwk"`

		Relationships []string `json:"relationships"`
	} `json:"verificationMethodKeys"`
}

type ServiceEndpoint struct {
	Nodes          []string `json:"nodes"`
	SigningKeys    []string `json:"signingKeys"`
	EncryptionKeys []string `json:"encryptionKeys"`
}

type didService struct {
	Id      string            `json:"id"`
	Type    string            `json:"type"`
	Service []ServiceEndpoint `json:"service"`
}

type didDocument struct {
	Id      string       `json:"id"`
	Service []didService `json:"service"`
	Context interface{}  `json:"@context"`
	// these keys may be either:
	// - a string, in which case it refers to another key in a different entry
	// - a map/object which containts the key in question
	KeyAgreement       []interface{} `json:"keyAgreement"`
	Authentication     []interface{} `json:"authentication"`
	VerificationMethod []interface{} `json:"verificationMethod"`
	// AlsoKnownAs is an array of URIs that represent alternative identifiers for the DID subject
	// https://www.w3.org/TR/did-1.0/#also-known-as
	AlsoKnownAs []string `json:"alsoKnownAs,omitempty"`
}
