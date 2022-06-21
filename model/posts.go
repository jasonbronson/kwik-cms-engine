package model

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Post struct {
	ID          string     `gorm:"primary_key;column:id;type:VARCHAR;size:255;unique;not null;" json:"id"`
	Title       string     `gorm:"column:title;type:VARCHAR;size:255;" json:"title"`
	Slug        string     `gorm:"column:slug;type:VARCHAR;size:255;unique;" json:"slug"`
	Content     string     `gorm:"column:content;type:TEXT;" json:"content"`
	Description string     `gorm:"column:description;type:VARCHAR;size:255;" json:"description"`
	UserID      string     `gorm:"column:user_id;type:VARCHAR;size:255;" json:"user_id"`
	User        User       `gorm:"foreignKey:ID;references:UserID" json:"user"`
	Status      string     `gorm:"column:status;type:VARCHAR;size:255;" json:"status"`
	Categories  []Category `gorm:"many2many:categories_post_links;ForeignKey:ID;References:ID"`
	Tags        []Tag      `gorm:"many2many:tags_post_links;ForeignKey:ID;References:ID"`
	MediaID     string     `gorm:"column:media_id;type:VARCHAR;size:255;" json:"media_id"`
	Media       Media      `gorm:"foreignKey:ID;references:MediaID" json:"media"`
	PublishDate time.Time  `gorm:"column:publish_date;type:TIMESTAMP;" json:"publish_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *Post) TableName() string {
	return "post"
}

func (u *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		newID, _ := uuid.NewV4()
		u.ID = uuid.Must(newID, nil).String()
	}
	return nil
}

func (a *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	tx.Where("post_id = ?", a.ID).Delete(CategoriesPostLinks{})
	tx.Where("post_id = ?", a.ID).Delete(TagsPostLinks{})
	return nil
}
