package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Category struct {
	ID        string    `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Name      string    `gorm:"column:name;type:VARCHAR;size:255;unique;" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `gorm:"many2many:categories_post_links;ForeignKey:ID;References:ID"`
}

func (u *Category) TableName() string {
	return "category"
}

func (u *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}
