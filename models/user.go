package models

import (
	"time"

	"github.com/rogery1999/go-gorm-rest-api/schemas"
)

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

func (u *User) MapUserToUserInfoResponse() *schemas.UserInfoResponse {
	return &schemas.UserInfoResponse{
		FirstName:  u.FirstName,
		MiddleName: u.MiddleName,
		LastName:   u.LastName,
		Email:      u.Email,
		Birthday:   u.Birthday.Format("2006-01-02"),
	}
}
