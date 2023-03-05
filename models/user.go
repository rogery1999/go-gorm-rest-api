package models

import "time"

type UserDTO struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Age  uint   `json:"age"`
}

type User struct {
	UUID       uint64 `gorm:"primaryKey;<-:false"`
	FirstName  string
	MiddleName string
	LastName   string
	Email      string `gorm:"uniqueIndex"`
	Password   string
	Birthday   time.Time `gorm:"index:,sort:asc"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
