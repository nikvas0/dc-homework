package objects

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	ID           uint32 `gorm:"primary_key"`
	Email        string `gorm:"unique_index:email"`
	PasswordHash string
	Confirmed    bool
}

type UserData struct {
	Email    string
	Password string
}

type Session struct {
	ID           uint32 `gorm:"primary_key"`
	UserID       uint32 `gorm:"index:sessionuserindex"`
	RefreshToken string `gorm:"index:sessionrefreshtoken"`
	ExpireAt     time.Time
}

type Token struct {
	UserID uint32
	Email  string
	jwt.StandardClaims
}

type ConfirmationToken struct {
	UserID uint32 `gorm:"primary_key"`
	Token  string `gorm:"unique_index:ctokenindex"`
	jwt.StandardClaims
}

type Notification struct {
	Email string
	Text  string
}
