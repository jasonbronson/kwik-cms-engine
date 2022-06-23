package model

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type DynamicField struct {
	ID           string         `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Label        string         `gorm:"column:label;type:VARCHAR;size:255;" json:"label"`
	Name         string         `gorm:"column:name;type:VARCHAR;size:255;" json:"name"`
	Type         string         `gorm:"column:type;type:VARCHAR;size:255;" json:"type"`
	Instructions string         `gorm:"column:instructions;type:VARCHAR;size:255;" json:"instructions"`
	Groups       []DynamicGroup `gorm:"many2many:dynamic_groups_fields;ForeignKey:ID;References:ID"`
}

func (f *DynamicField) TableName() string {
	return "dynamic_fields"
}

func (u *DynamicField) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}
