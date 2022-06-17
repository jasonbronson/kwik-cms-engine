package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID        string    `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Name      string    `gorm:"column:name;type:VARCHAR;size:255;unique;" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *Role) TableName() string {
	return "role"
}

func (u *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}
