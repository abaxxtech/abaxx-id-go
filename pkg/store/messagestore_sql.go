package store

import (
	"errors"
	"fmt"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/models"
	"github.com/fxamacker/cbor/v2"
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/multiformats/go-multihash"
	"gorm.io/gorm"
)

type GormMessageStore struct {
	db     *gorm.DB
	config MessageStoreSQLConfig
}

func NewMessageStoreSQL(config MessageStoreSQLConfig) (*GormMessageStore, error) {
	return &GormMessageStore{
		config: config,
	}, nil
}

func (mss *GormMessageStore) Open() error {
	db, err := models.GetDB(mss.config.DBConfig)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	mss.db = db
	return nil
}

func (mss *GormMessageStore) Close() error {
	// Don't close the actual connection since it's managed by singleton
	mss.db = nil
	return nil
}

func (mss *GormMessageStore) Get(tenant string, messageCid string) (interface{}, error) {
	var messageStore models.MessageStore
	if err := mss.db.Where("tenant = ? AND message_cid = ?", tenant, messageCid).First(&messageStore).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	// Decode the EncodedMessageBytes
	var message map[string]interface{}
	err := cbor.Unmarshal(messageStore.EncodedMessageBytes, &message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func getStringValue(m KeyValues, key string) string {
	if v, ok := m[key]; ok && v != nil {
		return v.(string)
	}
	return ""
}

func (mss *GormMessageStore) Put(tenant string, message GenericMessage, indexes KeyValues, options *MessageStoreOptions) error {
	encodedMessage, err := cbornode.WrapObject(message, multihash.SHA2_256, -1)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	messageCid := encodedMessage.Cid()

	messageStore := models.MessageStore{
		Tenant:               tenant,
		MessageCid:           messageCid.String(),
		EncodedMessageBytes:  encodedMessage.RawData(),
		EncodedData:          getStringValue(indexes, "encodedData"),
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
	}

	result := mss.db.Create(&messageStore)
	if result.Error != nil {
		return fmt.Errorf("failed to insert message: %w", result.Error)
	}

	return nil
}

func (mss *GormMessageStore) Query(tenant string, filters []DataFilter, messageSort *MessageSort, pagination *Pagination, options *MessageStoreOptions) ([]GenericMessage, string, error) {
	query := mss.db.Model(&models.MessageStore{}).Where("tenant = ?", tenant)

	// Apply filters
	for _, filter := range filters {
		query = query.Where(fmt.Sprintf("%s %s ?", filter.Property, filter.Operator), filter.Value)
	}

	// Apply sorting
	if messageSort != nil && messageSort.Property != "" {
		direction := messageSort.Direction
		if direction == "" {
			direction = "ASC"
		}
		query = query.Order(fmt.Sprintf("%s %s", messageSort.Property, direction))
	}

	// Apply pagination
	if pagination != nil && pagination.Limit > 0 {
		query = query.Limit(pagination.Limit)
		if pagination.Cursor != "" {
			offset, err := decodeCursor(pagination.Cursor)
			if err != nil {
				return nil, "", fmt.Errorf("invalid cursor: %w", err)
			}
			query = query.Offset(offset)
		}
	}

	var messages []models.MessageStore
	if err := query.Find(&messages).Error; err != nil {
		return nil, "", fmt.Errorf("failed to query messages: %w", err)
	}

	var genericMessages []GenericMessage
	for _, msg := range messages {
		var genericMessage GenericMessage
		if err := cbor.Unmarshal(msg.EncodedMessageBytes, &genericMessage); err != nil {
			return nil, "", fmt.Errorf("failed to unmarshal message: %w", err)
		}
		if msg.EncodedData != "" {
			if messageMap, ok := genericMessage.(map[string]interface{}); ok {
				messageMap["encodedData"] = msg.EncodedData
			}
		}
		genericMessages = append(genericMessages, genericMessage)
	}

	var nextCursor string
	if pagination != nil && len(messages) == pagination.Limit {
		offset := 0
		if pagination.Cursor != "" {
			offset, _ = decodeCursor(pagination.Cursor)
		}
		nextCursor = encodeCursor(offset + len(messages))
	}

	return genericMessages, nextCursor, nil
}

func (mss *GormMessageStore) Delete(tenant string, cidString string, options *MessageStoreOptions) error {
	result := mss.db.Where("tenant = ? AND message_cid = ?", tenant, cidString).Delete(&models.MessageStore{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete message: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("message not found")
	}
	return nil
}

func (mss *GormMessageStore) Clear() error {
	return mss.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.MessageStore{}).Error
}

// Helper functions for cursor encoding/decoding
func encodeCursor(offset int) string {
	return fmt.Sprintf("%d", offset)
}

func decodeCursor(cursor string) (int, error) {
	var offset int
	_, err := fmt.Sscanf(cursor, "%d", &offset)
	if err != nil {
		return 0, err
	}
	return offset, nil
}
