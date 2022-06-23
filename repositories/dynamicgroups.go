package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetDynamicGroups(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var dynamicgroups []model.DynamicGroup
	q := db.Model(&dynamicgroups)

	for i, param := range params.FilterBy {
		switch param {
		case "id":
			q.Scopes(ByID(params.FilterValue[i]))
		case "title":
			q.Scopes(FilterTitle(params.FilterValue[i]))

		}
	}

	params.Total = Count(db, &dynamicgroups)
	params.ResultTotal = Count(q, &dynamicgroups)
	q.Preload("DynamicFields").Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&dynamicgroups)
	return metaBuild(dynamicgroups, params)
}

func GetDynamicGroup(db *gorm.DB, id string) *response.Response {
	var dynamicgroup model.DynamicGroup
	q := db.Preload("DynamicFields").Where("id = ?", id)
	q.Find(&dynamicgroup)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(dynamicgroup, params)
}

func CreateDynamicGroup(db *gorm.DB, DynamicGroup model.DynamicGroup) error {
	return db.Create(&DynamicGroup).Error
}
func UpdateDynamicGroup(db *gorm.DB, DynamicGroup model.DynamicGroup) error {
	return db.Select("*").Updates(&DynamicGroup).Error
}

func DeleteDynamicGroup(db *gorm.DB, ID string) error {
	a := model.DynamicGroup{
		ID: ID,
	}
	return db.Delete(&a).Error
}
