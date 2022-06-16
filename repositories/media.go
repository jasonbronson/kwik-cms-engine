package repositories

import (
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	model "github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"gorm.io/gorm"
)

func GetMedia(db *gorm.DB, params helpers.DefaultParameters) *response.Response {
	var media []model.Media
	q := db.Model(&media)

	for i, param := range params.FilterBy {
		switch param {
		case "id":
			q.Scopes(ByID(params.FilterValue[i]))
		case "name":
			q.Scopes(FilterName(params.FilterValue[i]))
		}
	}

	q.Order(params.SortOrder).Limit(params.PageSize).Offset(params.PageOffset).Find(&media)
	return metaBuild(media, params)
}

func CreateMedia(db *gorm.DB, media model.Media) (string, error) {
	return media.ID, db.Create(&media).Error
}

func DeleteMedia(db *gorm.DB, mediaID string) error {
	return db.Delete(&model.Media{}, mediaID).Error
}
