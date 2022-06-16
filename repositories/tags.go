package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetTags(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var Tags []model.Tag
	q := db.Model(&Tags)

	if len(params.FilterBy) != 0 {
		switch params.FilterBy[0] {
		case "id":
			q.Scopes(ByID(params.FilterValue[0]))
		case "name":
			q.Scopes(FilterName(params.FilterValue[0]))
		}
	}

	params.Total = Count(db, &Tags)
	params.ResultTotal = Count(q, &Tags)
	q.Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&Tags)
	return metaBuild(Tags, params)
}
func GetTag(db *gorm.DB, TagID string) *response.Response {
	var Tag model.Tag
	db.Where("id = ?", TagID).First(&Tag)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(Tag, params)
}
func CreateTag(db *gorm.DB, Tag model.Tag) (string, error) {
	e := db.Create(&Tag).Error
	return Tag.ID, e
}
func UpdateTag(db *gorm.DB, Tag model.Tag) error {
	return db.Updates(&Tag).Error
}
func DeleteTag(db *gorm.DB, TagID string) error {
	Tag := &model.Tag{
		ID: TagID,
	}
	return db.Delete(&Tag).Error
}
