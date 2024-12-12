package models

import (
	"gorm.io/gorm"
)

// DataStore represents the main data storage table
type DataStore struct {
	gorm.Model
	Tenant      string `gorm:"not null;index:idx_tenant_data"`
	DataCid     string `gorm:"size:60;not null;index:idx_tenant_data"`
	EncodedData []byte `gorm:"type:bytea;not null"`
	DataFormat  string `gorm:"index"`
	Schema      string `gorm:"index"`
	DataSize    int64  `gorm:"index"`
	DateCreated string `gorm:"index"`
}
