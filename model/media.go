package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Media struct {
	ID        string    `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Name      string    `gorm:"column:name;type:VARCHAR;size:255;" json:"name"`
	URL       string    `gorm:"column:url;type:VARCHAR;size:255;" json:"url"`
	AltText   string    `gorm:"column:alt_text;type:VARCHAR;size:255;" json:"alt_text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (f *Media) TableName() string {
	return "media"
}

func (u *Media) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}
