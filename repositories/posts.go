package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetPosts(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var posts []model.Post
	q := db.Model(&posts)

	for i, param := range params.FilterBy {
		switch param {
		case "id":
			q.Scopes(ByID(params.FilterValue[i]))
		case "title":
			q.Scopes(FilterTitle(params.FilterValue[i]))
		case "description":
			q.Scopes(FilterDescription(params.FilterValue[i]))
		case "content":
			q.Scopes(FilterContent(params.FilterValue[i]))
		}
	}

	params.Total = Count(db, &posts)
	params.ResultTotal = Count(q, &posts)
	q.Preload("Author").Preload("Categories").Preload("Tags").Preload("Media").Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&posts)
	return metaBuild(posts, params)
}

func GetPostByID(db *gorm.DB, id string) *response.Response {
	var post model.Post
	q := db.Preload("Categories").Preload("Tags").Preload("Author").Preload("Media").Where("id = ?", id)
	q.Find(&post)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(post, params)
}

func CreatePost(db *gorm.DB, Post model.Post) error {
	return db.Create(&Post).Error
}
func UpdatePost(db *gorm.DB, Post model.Post) error {
	return db.Select("*").Updates(&Post).Error
}

func DeletePost(db *gorm.DB, PostID string) error {
	a := model.Post{
		ID: PostID,
	}
	return db.Delete(&a).Error
}
