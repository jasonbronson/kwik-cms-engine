package repositories

import (
	"strings"

	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var Users []model.User
	q := db.Model(&Users)

	if len(params.FilterBy) != 0 {
		switch params.FilterBy[0] {
		case "email":
			q.Scopes(FilterEmail(params.FilterValue[0]))
		}
	}

	params.Total = Count(db, &Users)
	params.ResultTotal = Count(q, &Users)
	q.Omit("Password").Preload("Role").Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&Users)
	return metaBuild(Users, params)
}
func GetUser(db *gorm.DB, UserID string) *response.Response {
	var User model.User
	db.Omit("Password").Preload("Role").Where("id = ?", UserID).First(&User)
	params := helpers.DefaultParameters{
		PageSize:   1,
		PageOffset: 0,
		Total:      1,
	}
	return metaBuild(User, params)
}
func GetUserByEmail(db *gorm.DB, email string) model.User {
	var User model.User
	email = strings.ToLower(email)
	db.Where("LOWER(email) = ?", email).First(&User)
	return User
}
func CreateUser(db *gorm.DB, User model.User) error {
	return db.Create(&User).Error
}
func UpdateUser(db *gorm.DB, User model.User) error {
	return db.Updates(&User).Error
}
func DeleteUser(db *gorm.DB, UserID string) error {
	User := &model.User{
		ID: UserID,
	}
	return db.Where("id = ?", UserID).Delete(&User).Error
}
