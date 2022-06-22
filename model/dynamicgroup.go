package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type DynamicGroup struct {
	ID        string         `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Title     string         `gorm:"column:title;type:VARCHAR;size:255;" json:"title"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	Fields    []DynamicField `gorm:"many2many:dynamic_groups_fields;ForeignKey:ID;References:ID"`
}

func (u *DynamicGroup) TableName() string {
	return "dynamic_groups"
}

func (u *DynamicGroup) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}

func (a *DynamicGroup) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Where("dynamic_group_id = ?", a.ID).Delete(DynamicGroupsFields{})
	return nil
}
