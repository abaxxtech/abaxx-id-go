package models

import (
	"fmt"
	"sync"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	instance *gorm.DB
	once     sync.Once
)

// GetDB returns a singleton instance of the database connection
func GetDB(config config.DBConfig) (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host,
			config.Port,
			config.User,
			config.Password,
			config.DBName,
			config.SSLMode,
		)

		instance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}

		// Auto-migrate all models in one place
		err = instance.AutoMigrate(
			&MessageStore{},
			&DataStore{},
			&DataStoreReference{},
			&EventLog{},
		)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return instance, nil
}
