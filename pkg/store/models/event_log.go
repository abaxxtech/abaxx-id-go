package models

import (
	"gorm.io/gorm"
)

// EventLog represents system events and audit logs
type EventLog struct {
	gorm.Model
	Tenant               string `gorm:"not null;index:idx_tenant_event"`
	EventType            string `gorm:"not null;index"`
	MessageCid           string `gorm:"size:60;index"`
	Interface            string `gorm:"index"`
	Method               string `gorm:"index"`
	Schema               string
	DataCid              string `gorm:"size:60;index"`
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
	Metadata             []byte `gorm:"type:jsonb"`
	Timestamp            string `gorm:"index"`
}

// TableName overrides the table name
func (EventLog) TableName() string {
	return "event_logs"
}

// Indexes to be created
func (EventLog) Indexes() []string {
	return []string{
		"idx_tenant_event",
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
