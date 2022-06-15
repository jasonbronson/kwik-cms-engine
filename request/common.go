package request

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
)

const (
	SortPublishDate = "publish"
	SortUpdateDate  = "updated_at"
	PageSize        = 20
)

func defaultParameters(g *gin.Context) helpers.DefaultParameters {

	page, err := strconv.Atoi(g.Request.URL.Query().Get("pagination[pageOffset]"))
	if err != nil {
		page = 0 //default to the first page if we don't get a valid offset
	}
	pageSize, err := strconv.Atoi(g.Request.URL.Query().Get("pagination[pageSize]"))
	if err != nil {
		pageSize = PageSize //default to pageSize items per page if a valid size isn't given
	}
	if pageSize == 0 || pageSize > 1000 {
		pageSize = PageSize
	}
	sortColumn := g.Request.URL.Query().Get("sort[0]")
	sortBy := g.Request.URL.Query().Get("sortby[0]")
	sortOrder := GetSortQueryString(sortColumn, sortBy)
	filters := GetFilterQueryString(g.Request)

	DefaultParameters := helpers.DefaultParameters{
		PageOffset:  page,
		PageSize:    pageSize,
		SortOrder:   sortOrder,
		FilterBy:    filters.FilterBy,
		FilterValue: filters.FilterValue,
		FilterQuery: filters.FilterQuery,
	}
	return DefaultParameters

}

func GetSortQueryString(sortColumn, sortBy string) string {
	var sortOrder string
	switch sortColumn {
	case "test":
		sortOrder = "test " + getSortBy(sortBy)
	default:
		sortOrder = "created_at " + getSortBy(sortBy)
	}
	return sortOrder
}
func GetFilterQueryString(r *http.Request) helpers.DefaultParameters {

	p := helpers.DefaultParameters{}
	for key, value := range r.URL.Query() {
		switch key {
		case "test":
			p.FilterBy = append(p.FilterBy, "test")
			p.FilterValue = append(p.FilterValue, value[0])
		}
	}
	return p
}
func getSortBy(sortBy string) string {
	switch sortBy {
	case "desc":
		sortBy = "desc"
	default:
		sortBy = "asc"
	}
	return sortBy
}
