package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Email     string `gorm:"column:email;type:VARCHAR;size:255;unique;" json:"email"`
	Password  string `gorm:"column:password;type:VARCHAR;size:255;" json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}
