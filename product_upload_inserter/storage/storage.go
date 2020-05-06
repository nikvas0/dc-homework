package storage

import (
	"errors"
	"log"
	"time"

	"product_upload_inserter/objects"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const connectRetries = 10
const sleepBetweenConnectRetriesDuration = 2 * time.Second

var db *gorm.DB

var errorNotFound = errors.New("not found")

func Init(database string, source string) error {
	if db != nil {
		return nil
	}

	var dbLocal *gorm.DB
	var err error
	counter := 0
	for {
		dbLocal, err = gorm.Open(database, source)
		if err != nil {
			counter++
			if counter == connectRetries {
				return err
			}
			log.Printf("Failed to connect to database: %v. Retrying...", err)
			time.Sleep(sleepBetweenConnectRetriesDuration)
		} else {
			break
		}
	}
	log.Println("Connected to database.")

	err = dbLocal.AutoMigrate(&objects.Product{}).Error
	if err != nil {
		return err
	}

	db = dbLocal

	db.LogMode(true)

	return nil
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func checkAndLogError(connection *gorm.DB) error {
	if connection.Error != nil {
		log.Printf("Storage error: %v.", connection.Error)
	}
	return connection.Error
}

func InsertProducts(products []objects.Product) error {
	tx := db.Begin()
	for _, product := range products {
		err := db.Exec("INSERT INTO \"products\" (\"id\",\"name\",\"category\") VALUES (?, ?, ?) ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, category=EXCLUDED.category;", product.ID, product.Name, product.Category).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
