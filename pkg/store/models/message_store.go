package models

import (
	"gorm.io/gorm"
)

// MessageStoreModel represents a DWN message in the database
type MessageStore struct {
	gorm.Model
	Tenant               string `gorm:"not null;index:idx_tenant_message"`
	MessageCid           string `gorm:"size:60;not null;index:idx_tenant_message"`
	EncodedMessageBytes  []byte `gorm:"type:bytea;not null"`
	EncodedData          string
	Interface            string `gorm:"index"`
	Method               string `gorm:"index"`
	Schema               string
	DataCid              string
	DataSize             string
	DateCreated          string `gorm:"index"`
	MessageTimestamp     string
	DataFormat           string
	IsLatestBaseState    string
	Published            string
	Author               string `gorm:"index"`
	RecordId             string `gorm:"index"`
	EntryId              string
	DatePublished        string
	Latest               string
	Protocol             string `gorm:"index"`
	DateExpires          string
	Description          string
	GrantedTo            string `gorm:"index"`
	GrantedBy            string `gorm:"index"`
	GrantedFor           string
	PermissionsRequestId string
	Attester             string
	ProtocolPath         string
	Recipient            string `gorm:"index"`
	ContextId            string
	ParentId             string `gorm:"index"`
	PermissionsGrantId   string
}

// TableName overrides the table name
func (MessageStore) TableName() string {
	return "message_store"
}

// Indexes to be created
func (MessageStore) Indexes() []string {
	return []string{
		"idx_tenant_message",
		"idx_interface",
		"idx_method",
		"idx_author",
		"idx_record_id",
		"idx_protocol",
		"idx_granted_to",
		"idx_granted_by",
		"idx_recipient",
		"idx_parent_id",
		"idx_date_created",
	}
}
