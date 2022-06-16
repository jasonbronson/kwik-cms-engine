package model

type TagsPostLinks struct {
	TagID  string `gorm:"column:tag_id;type:VARCHAR;size:255;" json:"tag_id"`
	PostID string `gorm:"column:post_id;type:VARCHAR;size:255;" json:"post_id"`
}

func (t *TagsPostLinks) TableName() string {
	return "tags_post_links"
}
