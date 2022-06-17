package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string    `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Email       string    `gorm:"column:email;type:VARCHAR;size:255;unique;" json:"email"`
	Password    string    `gorm:"column:password;type:VARCHAR;size:255;" json:"password"`
	RoleID      string    `gorm:"column:role_id;type:VARCHAR;size:255;" json:"role_id"`
	Role        Role      `gorm:"foreignKey:ID;references:RoleID" json:"role"`
	FirstName   string    `gorm:"column:first_name;type:VARCHAR;size:255;" json:"first_name"`
	LastName    string    `gorm:"column:last_name;type:VARCHAR;size:255;" json:"last_name"`
	Title       string    `gorm:"column:title;type:VARCHAR;size:255;" json:"title"`
	Content     string    `gorm:"column:content;type:VARCHAR;size:255;" json:"content"`
	FacebookUrl string    `gorm:"column:facebook_url;type:VARCHAR;size:255;" json:"facebook_url"`
	Instagram   string    `gorm:"column:instagram_url;type:VARCHAR;size:255;" json:"instagram_url"`
	Twitter     string    `gorm:"column:twitter_url;type:VARCHAR;size:255;" json:"twitter_url"`
	Youtube     string    `gorm:"column:youtube_url;type:VARCHAR;size:255;" json:"youtube_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
