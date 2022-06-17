package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetRoles(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var Roles []model.Role
	q := db.Model(&Roles)

	if len(params.FilterBy) != 0 {
		switch params.FilterBy[0] {
		case "id":
			q.Scopes(ByID(params.FilterValue[0]))
		case "name":
			q.Scopes(FilterName(params.FilterValue[0]))
		}
	}

	params.Total = Count(db, &Roles)
	params.ResultTotal = Count(q, &Roles)
	q.Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&Roles)
	return metaBuild(Roles, params)
}
func GetRole(db *gorm.DB, RoleID string) *response.Response {
	var Role model.Role
	db.Where("id = ?", RoleID).First(&Role)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(Role, params)
}
func CreateRole(db *gorm.DB, Role model.Role) (string, error) {
	e := db.Create(&Role).Error
	return Role.ID, e
}
func UpdateRole(db *gorm.DB, Role model.Role) error {
	return db.Updates(&Role).Error
}
func DeleteRole(db *gorm.DB, RoleID string) error {
	Role := &model.Role{
		ID: RoleID,
	}
	return db.Delete(&Role).Error
}
