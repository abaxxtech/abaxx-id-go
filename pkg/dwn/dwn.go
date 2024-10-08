package dwn

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/abaxxtech/abaxx-id-go/pkg/store"
)

type Signature struct {
	Protected string `json:"protected"`
	Signature string `json:"signature"`
}

type GeneralJws struct {
	Payload    string      `json:"payload"`
	Signatures []Signature `json:"signatures"`
}

type Authorization interface {
	signature() GeneralJws
}

type PlainAuthorization struct {
	Signature GeneralJws `json:"signature"`
}

func (a PlainAuthorization) signature() GeneralJws {
	return a.Signature
}

// Three types of signatures:
// plain -- signature only
// AuthOwner
// - Signature, AuthorDelegatedGrant, OwnerSignature, OwnerDelegatedGrant
// AuthDelegatedGrant
// - Signature, AuthorDelegatedGrant

type AuthorizationDelegatedGrant struct {
	Signature            GeneralJws      `json:"signature"`
	AuthorDelegatedGrant *DelegatedGrant `json:"authorDelegatedGrant"`
}

func (a *AuthorizationDelegatedGrant) signature() GeneralJws {
	return a.Signature
}

type AuthorizationOwner struct {
	Signature      GeneralJws `json:"signature"`
	OwnerSignature GeneralJws `json:"ownerSignature"`

	AuthorDelegatedGrant *DelegatedGrant `json:"authorDelegatedGrant,omitempty"`
	OwnerDelegatedGrant  *DelegatedGrant `json:"ownerDelegatedGrant,omitempty"`
}

func (a *AuthorizationOwner) signature() GeneralJws {
	return a.Signature
}

type DelegatedGrantDescriptor struct {
	RecordId    string `json:"recordId"`
	EncodedData string `json:"encodedData"`
}

type DelegatedGrant struct {
	// Authorization -> Signature --- ??
	Authorization PlainAuthorization       `json:"authorization"`
	Descriptor    DelegatedGrantDescriptor `json:"descriptor"`
}

type RecordsRead struct {
	Authorization AuthorizationDelegatedGrant `json:"authorization"`
	Descriptor    struct {
		Interface string
		Method    string
		// iso timestamp
		MessageTimestamp string
		// TODO make this properly typed.
		Filter map[string]interface{}
	} `json:"descriptor"`
}

type Descriptor struct {
	// Required
	Interface        string
	Method           string
	DataCid          DataCid
	DataSize         int64
	DateCreated      string
	MessageTimestamp string
	DataFormat       string

	Recipient     DID
	Protocol      string
	ProtocolPath  string
	Schema        string
	Tags          map[string]interface{}
	ParentId      MessageCid
	Published     bool
	DatePublished string
}

type RecordsWrite struct {
	// These three are required
	RecordId      string             `json:"recordId"`
	Authorization AuthorizationOwner `json:"authorization"`
	Descriptor    Descriptor         `json:"descriptor"`

	ContextId   string     `json:"contextId,omitempty"`
	Attestation GeneralJws `json:"attestation,omitempty"`
	Encryption  struct {
		Algorithm            string `json:"algorithm"`
		InitializationVector string `json:"initializationVector"`
		KeyEncryption        []struct {
			RootKeyId string
			// string-based enum of:
			// dataFormats | protocolContext | protocolPath | schemas
			DerivationScheme          string
			DerivedPublicKey          map[string]interface{}
			Algorithm                 string
			EncryptedKey              string
			InitializationVector      string
			EmphemeralPublicKey       map[string]interface{}
			MessageAuthenticationCode string
		} `json:"keyEncryption"`
	} `json:"encryption"`
}

type MessagesGet struct {
	Authorization PlainAuthorization `json:"authorization"`
	Descriptor    struct {
		Interface        string // const==Messages
		Method           string // const==Get
		MessageTimestamp string
		MessageCids      []string `json:"messageCids,omitempty"`
	} `json:"descriptor"`
}

func (m *AuthorizationDelegatedGrant) Author() string {

	return "none"
}

func (m *PlainAuthorization) Author() string {
	// if authorDelegatedGrant != nil
	//   Author = GetSigner(message.authorization.authorDelegatedGrant)
	// else
	//  Author = getSigner(m)

	// if m.Signature == nil {
	// 	return nil
	// }

	// TODO Utility functions for handling signatures

	// path:
	//  := m.Signature.Signatures[0]
	// extractDid( getKid ( checkThis ) )

	return "none"
}

type PermissionsGrant struct {
}

// TODO beef this up
type MessageHandler interface {
	Handle(dwn *Dwn) error
}

type MethodHandler interface {
	Handle(request *HandlerRequest) (UnionMessageReply, error)
}

type TenantGate interface {
	IsTenant(tenant string) (bool, error)
}

type AllowAllTenants struct{}

func (o AllowAllTenants) IsTenant(tenant string) (bool, error) {
	return true, nil
}

func NewAllowAllTenantGate() TenantGate {
	return &AllowAllTenants{}
}

type UnionMessageReply struct {
	Status Status
}

type Dwn struct {
	methodHandlers map[string]MethodHandler
	didResolver    *DidResolver
	messageStore   MessageStore
	dataStore      DataStore
	eventLog       EventLog
	tenantGate     TenantGate
	blockstore     *store.BlockstoreLevel
}

func NewDwn(config DwnConfig) (*Dwn, error) {
	if config.DidResolver == nil {
		config.DidResolver = NewDidResolver(nil, nil)
	}
	if config.TenantGate == nil {
		config.TenantGate = NewAllowAllTenantGate()
	}

	// Create a new BlockstoreLevel
	blockstore, err := store.NewBlockstoreLevel(config.BlockstoreLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to create blockstore: %w", err)
	}

	dwn := &Dwn{
		didResolver:    config.DidResolver,
		tenantGate:     config.TenantGate,
		messageStore:   config.MessageStore,
		dataStore:      config.DataStore,
		eventLog:       config.EventLog,
		blockstore:     blockstore,
		methodHandlers: map[string]MethodHandler{
			// "EventsGet":          NewEventsGetHandler(config.DidResolver, config.EventLog),
			// "EventsQuery":        NewEventsQueryHandler(config.DidResolver, config.EventLog),
			// "MessagesGet":        NewMessagesGetHandler(config.DidResolver, config.MessageStore, config.DataStore),
			// "PermissionsGrant":   NewPermissionsGrantHandler(config.DidResolver, config.MessageStore, config.EventLog),
			// "PermissionsRequest": NewPermissionsRequestHandler(config.DidResolver, config.MessageStore, config.EventLog),
			// "PermissionsRevoke":  NewPermissionsRevokeHandler(config.DidResolver, config.MessageStore, config.EventLog),
			// "ProtocolsConfigure": NewProtocolsConfigureHandler(config.DidResolver, config.MessageStore, config.DataStore, config.EventLog),
			// "ProtocolsQuery":     NewProtocolsQueryHandler(config.DidResolver, config.MessageStore, config.DataStore),
			// "RecordsDelete":      NewRecordsDeleteHandler(config.DidResolver, config.MessageStore, config.DataStore, config.EventLog),
			// "RecordsQuery":       NewRecordsQueryHandler(config.DidResolver, config.MessageStore, config.DataStore),
			// "RecordsRead":        NewRecordsReadHandler(config.DidResolver, config.MessageStore, config.DataStore),
			// "RecordsWrite":       NewRecordsWriteHandler(config.DidResolver, config.MessageStore, config.DataStore, config.EventLog),
		},
	}

	if err := dwn.Open(); err != nil {
		return nil, err
	}

	return dwn, nil
}

func (d *Dwn) Open() error {
	if err := d.messageStore.Open(); err != nil {
		return err
	}
	if err := d.dataStore.Open(); err != nil {
		return err
	}
	if err := d.eventLog.Open(); err != nil {
		return err
	}
	if err := d.blockstore.Open(); err != nil {
		return err
	}
	return nil
}

func (d *Dwn) Close() error {
	if err := d.messageStore.Close(); err != nil {
		return err
	}
	if err := d.dataStore.Close(); err != nil {
		return err
	}
	if err := d.eventLog.Close(); err != nil {
		return err
	}
	if err := d.blockstore.Close(); err != nil {
		return err
	}
	return nil
}

// This function steps thru the json document, following the keys in `paths`, returning the
// string, or empty if not found.
func getPathedStrNoErr(json map[string]interface{}, paths ...string) string {
	var current interface{} = json

	for _, path := range paths {
		if m, ok := current.(map[string]interface{}); ok {
			current, ok = m[path]
			if !ok {
				return ""
			}
		} else {
			return ""
		}
	}

	if str, ok := current.(string); ok {
		return str
	}

	return ""
}

func (d *Dwn) ProcessMessage(tenant string, rawMessage map[string]interface{}, dataStream io.Reader) (UnionMessageReply, error) {
	if err := d.validateTenant(tenant); err != nil {
		return UnionMessageReply{Status: Status{Code: 401, Detail: err.Error()}}, nil
	}

	if err := d.validateMessageIntegrity(rawMessage); err != nil {
		return UnionMessageReply{Status: Status{Code: 400, Detail: err.Error()}}, nil
	}

	//
	handlerKey := getPathedStrNoErr(rawMessage, "Descriptor", "Interface") + getPathedStrNoErr(rawMessage, "Descriptor", "Method")
	methodHandler, exists := d.methodHandlers[handlerKey]
	if !exists {
		return UnionMessageReply{}, errors.New("handler not found")
	}

	return methodHandler.Handle(&HandlerRequest{
		Tenant:     tenant,
		Message:    rawMessage,
		DataStream: dataStream,
	})
}

func (d *Dwn) validateTenant(tenant string) error {
	isTenant, err := d.tenantGate.IsTenant(tenant)
	if err != nil {
		return err
	}
	if !isTenant {
		return errors.New(tenant + " is not a tenant")
	}
	return nil
}

func (d *Dwn) validateMessageIntegrity(rawMessage map[string]interface{}) error {
	if getPathedStrNoErr(rawMessage, "Descriptor", "Interface") == "" ||
		getPathedStrNoErr(rawMessage, "Descriptor", "Method") == "" {
		return errors.New("both interface and method must be present")
	}

	// if err := ValidateJsonSchema(rawMessage); err != nil {
	// 	return err
	// }

	return nil
}

func (dwn *Dwn) authenticate(auth Authorization) error {
	_, err := json.Marshal(auth.signature())
	if err != nil {
		return err
	}

	// sig, err := jose.ParseSigned(string(reser),
	// 	[]jose.SignatureAlgorithm{jose.EdDSA, jose.HS256, jose.HS384})
	//if err != nil {
	//return err
	//}
	//fmt.Println("sig", sig)
	// for i, s := range sig.Signatures {
	// 	fmt.Println(i, "keyid", s.Header.KeyID)
	// 	// didResult, err := dwn.didResolver.Resolve(s.Header.KeyID)
	// 	// if err != nil {
	// 	// 	return err
	// 	// }
	// 	//fmt.Println(i, "result", didResult)
	// }

	// TODO implement authentication logic
	return nil
}

func (dwn *Dwn) authorize(tenant Tenant, auth Authorization) error {
	return nil
}
