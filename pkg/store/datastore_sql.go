package store

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/abaxxtech/abaxx-id-go/pkg/store/models"
	"gorm.io/gorm"
)

type GetGormResult struct {
	DataCid    string
	DataSize   int64
	DataStream io.Reader
}

type PutGormResult struct {
	DataCid  string
	DataSize int64
}

type AssociateGormResult struct {
	DataCid  string
	DataSize int64
}

type DataStoreSQL struct {
	db     *gorm.DB
	config MessageStoreSQLConfig
}

func NewDataStoreSQL(config MessageStoreSQLConfig) (*DataStoreSQL, error) {
	return &DataStoreSQL{
		config: config,
	}, nil
}

func (dss *DataStoreSQL) Open() error {
	db, err := models.GetDB(dss.config.DBConfig)
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Auto-migrate tables
	err = db.AutoMigrate(&models.DataStore{}, &models.DataStoreReference{})
	if err != nil {
		return fmt.Errorf("failed to create database schema: %w", err)
	}

	dss.db = db
	return nil
}

func (dss *DataStoreSQL) Close() error {
	dss.db = nil
	return nil
}

func (dss *DataStoreSQL) Get(tenant string, messageCid string, dataCid string) (string, int64, io.Reader, error) {
	if dss.db == nil {
		return "", 0, nil, fmt.Errorf("connection to database not open. Call `open` before using `get`")
	}

	var ref models.DataStoreReference
	result := dss.db.Where(&models.DataStoreReference{
		Tenant:     tenant,
		MessageCid: messageCid,
		DataCid:    dataCid,
	}).First(&ref)

	if result.Error == gorm.ErrRecordNotFound {
		return "", 0, nil, nil
	}
	if result.Error != nil {
		return "", 0, nil, fmt.Errorf("failed to check reference: %w", result.Error)
	}

	var data models.DataStore
	result = dss.db.Where(&models.DataStore{
		Tenant:  tenant,
		DataCid: dataCid,
	}).First(&data)

	if result.Error == gorm.ErrRecordNotFound {
		return "", 0, nil, nil
	}
	if result.Error != nil {
		return "", 0, nil, fmt.Errorf("failed to get data: %w", result.Error)
	}

	return dataCid, int64(len(data.EncodedData)), io.NopCloser(bytes.NewReader(data.EncodedData)), nil
}

func (dss *DataStoreSQL) Put(tenant string, messageCid string, dataCid string,
	dataStream io.Reader) (*PutGormResult, error) {
	if dss.db == nil {
		return nil, fmt.Errorf("connection to database not open. Call `open` before using `put`")
	}

	data, err := ioutil.ReadAll(dataStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read data stream: %w", err)
	}

	err = dss.db.Transaction(func(tx *gorm.DB) error {
		dataStore := models.DataStore{
			Tenant:      tenant,
			DataCid:     dataCid,
			EncodedData: data,
		}
		if err := tx.Create(&dataStore).Error; err != nil {
			return fmt.Errorf("failed to store data: %w", err)
		}

		ref := models.DataStoreReference{
			Tenant:     tenant,
			MessageCid: messageCid,
			DataCid:    dataCid,
		}
		if err := tx.Create(&ref).Error; err != nil {
			return fmt.Errorf("failed to create reference: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &PutGormResult{
		DataCid:  dataCid,
		DataSize: int64(len(data)),
	}, nil
}

func (dss *DataStoreSQL) Associate(tenant string, messageCid string, dataCid string) (*AssociateGormResult, error) {
	if dss.db == nil {
		return nil, fmt.Errorf("connection to database not open. Call `open` before using `associate`")
	}

	var data models.DataStore
	result := dss.db.Where(&models.DataStore{
		Tenant:  tenant,
		DataCid: dataCid,
	}).First(&data)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to check data: %w", result.Error)
	}

	// Create reference in a transaction to ensure atomicity
	err := dss.db.Transaction(func(tx *gorm.DB) error {
		ref := models.DataStoreReference{
			Tenant:     tenant,
			MessageCid: messageCid,
			DataCid:    dataCid,
		}

		// Check if reference already exists
		var existingRef models.DataStoreReference
		result := tx.Where(&ref).First(&existingRef)
		if result.Error == gorm.ErrRecordNotFound {
			// Create new reference if it doesn't exist
			if err := tx.Create(&ref).Error; err != nil {
				return fmt.Errorf("failed to create reference: %w", err)
			}
		} else if result.Error != nil {
			return fmt.Errorf("failed to check existing reference: %w", result.Error)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AssociateGormResult{
		DataCid:  dataCid,
		DataSize: int64(len(data.EncodedData)),
	}, nil
}

func (dss *DataStoreSQL) Delete(tenant string, messageCid string, dataCid string) error {
	if dss.db == nil {
		return fmt.Errorf("connection to database not open. Call `open` before using `delete`")
	}

	return dss.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(&models.DataStoreReference{
			Tenant:     tenant,
			MessageCid: messageCid,
			DataCid:    dataCid,
		}).Delete(&models.DataStoreReference{}).Error; err != nil {
			return fmt.Errorf("failed to delete reference: %w", err)
		}

		var count int64
		if err := tx.Model(&models.DataStoreReference{}).Where(&models.DataStoreReference{
			Tenant:  tenant,
			DataCid: dataCid,
		}).Count(&count).Error; err != nil {
			return fmt.Errorf("failed to count references: %w", err)
		}

		if count == 0 {
			if err := tx.Where(&models.DataStore{
				Tenant:  tenant,
				DataCid: dataCid,
			}).Delete(&models.DataStore{}).Error; err != nil {
				return fmt.Errorf("failed to delete data: %w", err)
			}
		}

		return nil
	})
}

func (dss *DataStoreSQL) Clear() error {
	if dss.db == nil {
		return fmt.Errorf("connection to database not open. Call `open` before using `clear`")
	}

	return dss.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.DataStoreReference{}).Error; err != nil {
			return err
		}
		return tx.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.DataStore{}).Error
	})
}
