package storage

import (
	"errors"
	"log"
	"time"

	"github.com/nikvas0/dc-homework/server/objects"
	"github.com/nikvas0/gorm"
	_ "github.com/nikvas0/gorm/dialects/postgres"
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

func IsNotFoundError(err error) bool {
	return err == errorNotFound
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

func CreateProduct(product *objects.Product) error {
	return checkAndLogError(db.Create(product))
}

func GetProductById(product *objects.Product, id uint32) error {
	return checkAndLogError(db.First(product, id))
}

func GetProductsPage(products *[]objects.Product, offset uint32, limit uint32) error {
	err := db.Order("id").Offset(offset).Limit(limit).Find(products).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return errorNotFound
	} else if err != nil {
		log.Printf("Storage error: %v.", err)
	}
	return err
}

func GetProductsCount() (uint32, error) {
	var count uint32
	err := db.Model(objects.Product{}).Count(&count).Error
	return count, err
}

func GetAllProducts(products *[]objects.Product) error {
	err := db.Order("id").Find(products).Error
	if err != nil && gorm.IsRecordNotFoundError(err) {
		return errorNotFound
	} else if err != nil {
		log.Printf("Storage error: %v.", err)
	}
	return err
}

func DeleteProductById(id uint32) error {
	return checkAndLogError(db.Delete(&objects.Product{ID: id}))
}

func UpdateProduct(product *objects.Product) error {
	return checkAndLogError(db.Model(product).Updates(product))
}

func checkAndLogError(connection *gorm.DB) error {
	if (connection.Error != nil && gorm.IsRecordNotFoundError(connection.Error)) || (connection.Error == nil && connection.RowsAffected == 0) {
		return errorNotFound
	} else if connection.Error != nil {
		log.Printf("Storage error: %v.", connection.Error)
	}
	return connection.Error
}
