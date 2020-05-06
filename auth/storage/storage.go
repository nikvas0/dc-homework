package storage

import (
	"errors"
	"log"
	"time"

	"auth/objects"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

const connectRetries = 10
const sleepBetweenConnectRetriesDuration = 2 * time.Second

var db *gorm.DB

var errorNotFound = errors.New("not found")
var errorAlreadyExists = errors.New("already exists")

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

	err = dbLocal.AutoMigrate(&objects.User{}).Error
	if err != nil {
		return err
	}

	err = dbLocal.AutoMigrate(&objects.Session{}).Error
	if err != nil {
		return err
	}

	err = dbLocal.AutoMigrate(&objects.ConfirmationToken{}).Error
	if err != nil {
		return err
	}

	db = dbLocal

	db.LogMode(true)

	user := objects.User{}
	user.Email = "admin"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("saltadmin"), bcrypt.DefaultCost)
	user.PasswordHash = string(hashedPassword)
	user.Confirmed = true
	user.Role = objects.AdminRole
	CreateUser(&user) // ignore error

	return nil
}

func IsErrorAlreadyExists(err error) bool {
	return err == errorAlreadyExists
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

func CreateUser(user *objects.User) error {
	connection := db.Where("email = ?", user.Email).First(user)
	if connection.Error != nil {
		log.Printf("Storage error: %v.", connection.Error)
	} else if connection.RowsAffected != 0 && !user.Confirmed {
		return errorAlreadyExists
	}
	connection = db.Create(user)
	if connection.Error != nil {
		log.Printf("Storage error: %v.", connection.Error)
	}
	return connection.Error
}

func CreateConfirmationToken(user *objects.User, token string) error {
	return checkAndLogError(db.Create(&objects.ConfirmationToken{UserID: user.ID, Token: token}))
}

func GetConfirmationToken(user *objects.User, token *string) error {
	obj := objects.ConfirmationToken{}
	err := checkAndLogError(db.First(obj, user.ID))
	token = &obj.Token
	return err
}

func DeleteConfirmationToken(user *objects.User) error {
	return checkAndLogError(db.Delete(&objects.ConfirmationToken{UserID: user.ID}))
}

func GetUserByEmail(email string, user *objects.User) error {
	return checkAndLogError(db.Where("email = ?", email).First(user))
}

func GetUserByID(id uint32, user *objects.User) error {
	return checkAndLogError(db.First(user, id))
}

func GetUserIDByConfirmationToken(token string, id *uint32) error {
	res := objects.ConfirmationToken{}
	err := checkAndLogError(db.Where("token = ?", token).First(&res))
	*id = res.UserID
	return err
}

func UpdateUserRoleByID(id uint32, role uint32) error {
	return checkAndLogError(db.Model(&objects.User{}).Where("id = ?", id).Update("role", role))
}

func SetUserConfirmedByID(id uint32) error {
	return checkAndLogError(db.Model(&objects.User{}).Where("id = ?", id).Update("confirmed", true))
}

func CreateSession(session *objects.Session) error {
	return checkAndLogError(db.Create(session))
}

func GetSessionByToken(token string, session *objects.Session) error {
	return checkAndLogError(db.Where("refresh_token = ? AND expire_at > ?", token, time.Now()).First(session))
}

func UpdateSession(session *objects.Session) error {
	return checkAndLogError(db.Save(session))
}

func checkAndLogError(connection *gorm.DB) error {
	if (connection.Error != nil && gorm.IsRecordNotFoundError(connection.Error)) || (connection.Error == nil && connection.RowsAffected == 0) {
		return errorNotFound
	} else if connection.Error != nil {
		log.Printf("Storage error: %v.", connection.Error)
	}
	return connection.Error
}
