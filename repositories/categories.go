package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetCategories(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var Categories []model.Category
	q := db.Model(&Categories)

	for i, param := range params.FilterBy {
		switch param {
		case "id":
			q.Scopes(ByID(params.FilterValue[i]))
		case "name":
			q.Scopes(FilterName(params.FilterValue[i]))
		}
	}

	params.Total = Count(db, &Categories)
	params.ResultTotal = Count(q, &Categories)
	q.Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&Categories)
	return metaBuild(Categories, params)
}
func GetCategory(db *gorm.DB, CategoryID string) *response.Response {
	var Category model.Category
	db.Where("id = ?", CategoryID).First(&Category)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(Category, params)
}
func CreateCategory(db *gorm.DB, Category model.Category) (string, error) {
	e := db.Create(&Category).Error
	return Category.ID, e
}
func UpdateCategory(db *gorm.DB, Category model.Category) error {
	return db.Updates(&Category).Error
}
func DeleteCategory(db *gorm.DB, CategoryID string) error {
	Category := &model.Category{}
	return db.Where("id = ?", CategoryID).Delete(&Category).Error
}
