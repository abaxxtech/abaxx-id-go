package dwn

import (
	"errors"
	"io"
)


type MethodHandler interface {
	Handle(request *HandlerRequest) (UnionMessageReply, error)
}

type TenantGate interface {
	IsTenant(tenant string) (bool, error)
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
}

func NewDwn(config DwnConfig) (*Dwn, error) {
	if config.DidResolver == nil {
		config.DidResolver = NewDidResolver()
	}
	if config.TenantGate == nil {
		config.TenantGate = NewAllowAllTenantGate()
	}

	dwn := &Dwn{
		didResolver:  config.DidResolver,
		tenantGate:   config.TenantGate,
		messageStore: config.MessageStore,
		dataStore:    config.DataStore,
		eventLog:     config.EventLog,
		methodHandlers: map[string]MethodHandler{
			"EventsGet":          NewEventsGetHandler(config.DidResolver, config.EventLog),
			"EventsQuery":        NewEventsQueryHandler(config.DidResolver, config.EventLog),
			"MessagesGet":        NewMessagesGetHandler(config.DidResolver, config.MessageStore, config.DataStore),
			"PermissionsGrant":   NewPermissionsGrantHandler(config.DidResolver, config.MessageStore, config.EventLog),
			"PermissionsRequest": NewPermissionsRequestHandler(config.DidResolver, config.MessageStore, config.EventLog),
			"PermissionsRevoke":  NewPermissionsRevokeHandler(config.DidResolver, config.MessageStore, config.EventLog),
			"ProtocolsConfigure": NewProtocolsConfigureHandler(config.DidResolver, config.MessageStore, config.DataStore, config.EventLog),
			"ProtocolsQuery":     NewProtocolsQueryHandler(config.DidResolver, config.MessageStore, config.DataStore),
			"RecordsDelete":      NewRecordsDeleteHandler(config.DidResolver, config.MessageStore, config.DataStore, config.EventLog),
			"RecordsQuery":       NewRecordsQueryHandler(config.DidResolver, config.MessageStore, config.DataStore),
			"RecordsRead":        NewRecordsReadHandler(config.DidResolver, config.MessageStore, config.DataStore),
			"RecordsWrite":       NewRecordsWriteHandler(config.DidResolver, config.MessageStore, config.DataStore, config.EventLog),
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
	return nil
}

func (d *Dwn) ProcessMessage(tenant string, rawMessage GenericMessage, dataStream io.Reader) (UnionMessageReply, error) {
	if err := d.validateTenant(tenant); err != nil {
		return UnionMessageReply{Status: Status{Code: 401, Detail: err.Error()}}, nil
	}

	if err := d.validateMessageIntegrity(rawMessage); err != nil {
		return UnionMessageReply{Status: Status{Code: 400, Detail: err.Error()}}, nil
	}

	handlerKey := rawMessage.Descriptor.Interface + rawMessage.Descriptor.Method
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

func (d *Dwn) validateMessageIntegrity(rawMessage GenericMessage) error {
	if rawMessage.Descriptor.Interface == "" || rawMessage.Descriptor.Method == "" {
		return errors.New("both interface and method must be present")
	}

	if err := ValidateJsonSchema(rawMessage); err != nil {
		return err
	}

	return nil
}
