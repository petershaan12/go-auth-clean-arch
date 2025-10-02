package model

type Role struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"type:varchar(50);unique"`
}
