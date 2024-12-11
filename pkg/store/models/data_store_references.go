package models

import (
	"gorm.io/gorm"
)

// DataStoreReference represents references between data objects
type DataStoreReference struct {
	gorm.Model
	Tenant        string `gorm:"not null;index:idx_tenant_ref"`
	SourceCid     string `gorm:"size:60;not null;index:idx_source"`
	TargetCid     string `gorm:"size:60;not null;index:idx_target"`
	ReferenceType string `gorm:"index"`
	DateCreated   string `gorm:"index"`
}
