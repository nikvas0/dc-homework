package objects

type Product struct {
	ID       uint32 `gorm:"primary_key"`
	Name     string
	Category uint32
}
