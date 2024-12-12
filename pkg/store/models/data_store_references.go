package models

import (
	"gorm.io/gorm"
)

// DataStoreReference represents references between data objects
type DataStoreReference struct {
	gorm.Model
	Tenant      string `gorm:"not null;index:idx_tenant_ref"`
	DataCid     string `gorm:"size:60;not null;index:idx_data_cid"`
	MessageCid  string `gorm:"size:60;not null;index:idx_message_cid"`
	DateCreated string `gorm:"index"`
}
