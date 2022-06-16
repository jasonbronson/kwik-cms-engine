package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetAuthors(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var Authors []model.Author
	q := db.Model(&Authors)

	if len(params.FilterBy) != 0 {
		switch params.FilterBy[0] {
		case "id":
			q.Scopes(ByID(params.FilterValue[0]))
		case "username":
			q.Scopes(FilterUsername(params.FilterValue[0]))
		case "email":
			q.Scopes(FilterEmail(params.FilterValue[0]))
		}
	}

	params.Total = Count(db, &Authors)
	params.ResultTotal = Count(q, &Authors)
	q.Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&Authors)
	return metaBuild(Authors, params)
}
func GetAuthorByID(db *gorm.DB, ID string) *response.Response {
	var Author model.Author
	db.Where("id = ?", ID).Preload("Media").First(&Author)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(Author, params)
}
func CreateAuthor(db *gorm.DB, Author model.Author) error {
	return db.Create(&Author).Error
}
func UpdateAuthor(db *gorm.DB, Author model.Author) error {
	return db.Updates(&Author).Error
}
func DeleteAuthor(db *gorm.DB, AuthorID string) error {
	Author := &model.Author{
		ID: AuthorID,
	}
	return db.Where("id = ?", AuthorID).Delete(&Author).Error
}
