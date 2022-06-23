package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	q.Preload("Fields").Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&dynamicgroups)
	return metaBuild(dynamicgroups, params)
}

func GetDynamicGroup(db *gorm.DB, id string) *response.Response {
	var dynamicgroup model.DynamicGroup
	q := db.Preload("Fields").Where("id = ?", id)
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
	tx := db.Begin()
	err := db.Omit(clause.Associations).Save(&DynamicGroup).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var fieldIds []string
	for _, field := range DynamicGroup.Fields {
		fieldIds = append(fieldIds, field.ID)
	}

	if len(fieldIds) == 0 {
		db.Where("dynamic_group_id=?", DynamicGroup.ID).Delete(&model.DynamicGroupsFields{})
	} else {
		db.Where("dynamic_field_id NOT IN ?", fieldIds).Where("dynamic_group_id=?", DynamicGroup.ID).Delete(&model.DynamicGroupsFields{})
	}

	for _, field := range DynamicGroup.Fields {
		if err := db.Save(&field).Error; err != nil {
			tx.Rollback()
			return err
		}
		link := model.DynamicGroupsFields{}
		if !contains(fieldIds, field.ID) {
			link.DynamicFieldID = field.ID
			link.DynamicGroupID = DynamicGroup.ID
			if err := db.Create(&link).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func DeleteDynamicGroup(db *gorm.DB, ID string) error {
	a := model.DynamicGroup{
		ID: ID,
	}
	return db.Delete(&a).Error
}
