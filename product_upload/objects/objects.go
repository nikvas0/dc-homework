package objects

import (
	"errors"
	"strings"
)

type Product struct {
	ID       uint32 `gorm:"primary_key"`
	Name     string
	Category uint32
}

func FixProduct(product *Product) error {
	product.Name = strings.TrimSpace(product.Name)
	if len(product.Name) == 0 {
		return errors.New("bad product name")
	}

	if product.Category == 0 {
		return errors.New("bad product category")
	}

	return nil
}

type UserData struct {
	ID    uint32
	Email string
}
