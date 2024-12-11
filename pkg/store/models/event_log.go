package models

import (
	"gorm.io/gorm"
)

// EventLog represents system events and audit logs
type EventLog struct {
	gorm.Model
	Tenant      string `gorm:"not null;index:idx_tenant_event"`
	EventType   string `gorm:"not null;index"`
	MessageCid  string `gorm:"size:60;index"`
	DataCid     string `gorm:"size:60;index"`
	Description string
	Metadata    []byte `gorm:"type:jsonb"`
	Timestamp   string `gorm:"index"`
}
