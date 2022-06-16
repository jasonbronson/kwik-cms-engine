package model

type CategoriesPostLinks struct {
	CategoryID string `gorm:"column:category_id;type:VARCHAR;size:255;" json:"category_id"`
	PostID     string `gorm:"column:post_id;type:VARCHAR;size:255;" json:"post_id"`
}

func (c *CategoriesPostLinks) TableName() string {
	return "categories_post_links"
}
