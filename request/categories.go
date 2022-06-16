package request

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/repositories"

	"github.com/jasonbronson/kwik-cms-engine/request/response"
	_ "github.com/jasonbronson/kwik-cms-engine/request/response"
)

func GetCategories(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetCategories(db, params)
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetCategory(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("categoryid")
	r := repositories.GetCategory(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutCategories(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var category model.Category
	if e := g.ShouldBindJSON(&category); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	e := repositories.UpdateCategory(db, category)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func SetCategories(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var category model.Category
	if e := g.ShouldBindJSON(&category); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	id, e := repositories.CreateCategory(db, category)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, fmt.Sprint(id))
}

func DeleteCategories(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	categoryID := g.Param("categoryid")
	e := repositories.DeleteCategory(db, categoryID)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
