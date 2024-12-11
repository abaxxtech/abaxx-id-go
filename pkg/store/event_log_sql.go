package store

import (
	"fmt"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/models"
	"gorm.io/gorm"
)

type EventLogSQL struct {
	db     *gorm.DB
	config MessageStoreSQLConfig
}

func NewEventLogSQL(config MessageStoreSQLConfig) (*EventLogSQL, error) {
	return &EventLogSQL{
		config: config,
	}, nil
}

func (els *EventLogSQL) Open() error {
	db, err := models.GetDB(els.config.DBConfig)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	els.db = db
	return nil
}

func (els *EventLogSQL) Close() error {
	els.db = nil
	return nil
}

func (els *EventLogSQL) Append(tenant string, messageCid string, indexes KeyValues) error {
	if els.db == nil {
		return fmt.Errorf("database connection not open")
	}

	eventLog := models.EventLog{
		Tenant:               tenant,
		MessageCid:           messageCid,
		EventType:            getStringValue(indexes, "eventType"),
		Interface:            getStringValue(indexes, "interface"),
		Method:               getStringValue(indexes, "method"),
		Schema:               getStringValue(indexes, "schema"),
		DataCid:              getStringValue(indexes, "dataCid"),
		DataSize:             getStringValue(indexes, "dataSize"),
		DateCreated:          getStringValue(indexes, "dateCreated"),
		MessageTimestamp:     getStringValue(indexes, "messageTimestamp"),
		DataFormat:           getStringValue(indexes, "dataFormat"),
		IsLatestBaseState:    getStringValue(indexes, "isLatestBaseState"),
		Published:            getStringValue(indexes, "published"),
		Author:               getStringValue(indexes, "author"),
		RecordId:             getStringValue(indexes, "recordId"),
		EntryId:              getStringValue(indexes, "entryId"),
		DatePublished:        getStringValue(indexes, "datePublished"),
		Latest:               getStringValue(indexes, "latest"),
		Protocol:             getStringValue(indexes, "protocol"),
		DateExpires:          getStringValue(indexes, "dateExpires"),
		Description:          getStringValue(indexes, "description"),
		GrantedTo:            getStringValue(indexes, "grantedTo"),
		GrantedBy:            getStringValue(indexes, "grantedBy"),
		GrantedFor:           getStringValue(indexes, "grantedFor"),
		PermissionsRequestId: getStringValue(indexes, "permissionsRequestId"),
		Attester:             getStringValue(indexes, "attester"),
		ProtocolPath:         getStringValue(indexes, "protocolPath"),
		Recipient:            getStringValue(indexes, "recipient"),
		ContextId:            getStringValue(indexes, "contextId"),
		ParentId:             getStringValue(indexes, "parentId"),
		PermissionsGrantId:   getStringValue(indexes, "permissionsGrantId"),
		Metadata:             []byte(getStringValue(indexes, "metadata")),
		Timestamp:            getStringValue(indexes, "timestamp"),
	}

	return els.db.Create(&eventLog).Error
}

func (els *EventLogSQL) GetEvents(tenant string, options *EventOptions) ([]string, error) {
	return els.QueryEvents(tenant, nil, options)
}

func (els *EventLogSQL) QueryEvents(tenant string, filters []DataFilter, options *EventOptions) ([]string, error) {
	if els.db == nil {
		return nil, fmt.Errorf("database connection not open")
	}

	query := els.db.Model(&models.EventLog{}).Where("tenant = ?", tenant)

	// Apply filters
	for _, filter := range filters {
		query = query.Where(fmt.Sprintf("%s %s ?", filter.Property, filter.Operator), filter.Value)
	}

	// Apply cursor-based pagination if cursor is provided
	if options != nil && options.Cursor != "" {
		var cursorEvent models.EventLog
		if err := els.db.Where("tenant = ? AND message_cid = ?", tenant, options.Cursor).First(&cursorEvent).Error; err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		query = query.Where("id > ?", cursorEvent.ID)
	}

	// Order by ID (watermark) ascending
	query = query.Order("id asc")

	var events []models.EventLog
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	// Extract message CIDs
	messageCids := make([]string, len(events))
	for i, event := range events {
		messageCids[i] = event.MessageCid
	}

	return messageCids, nil
}

func (els *EventLogSQL) DeleteEventsByCid(tenant string, messageCids []string) error {
	if els.db == nil {
		return fmt.Errorf("database connection not open")
	}

	if len(messageCids) == 0 {
		return nil
	}

	return els.db.Where("tenant = ? AND message_cid IN ?", tenant, messageCids).Delete(&models.EventLog{}).Error
}

func (els *EventLogSQL) Clear() error {
	if els.db == nil {
		return fmt.Errorf("database connection not open")
	}

	return els.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.EventLog{}).Error
}

type EventOptions struct {
	Cursor string
}
