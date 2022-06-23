package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/repositories"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
)

func GetDynamicGroups(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetDynamicGroups(db, params)
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetDynamicGroup(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("postid")
	var r *response.Response
	r = repositories.GetDynamicGroup(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutDynamicGroup(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	var dynamicgroup model.DynamicGroup
	if e := g.ShouldBindJSON(&dynamicgroup); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	e := repositories.UpdateDynamicGroup(db, dynamicgroup)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func PostDynamicGroup(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	var dynamicgroup model.DynamicGroup
	if e := g.ShouldBindJSON(&dynamicgroup); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	e := repositories.CreateDynamicGroup(db, dynamicgroup)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func DeleteDynamicGroup(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("id")
	e := repositories.DeleteDynamicGroup(db, id)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
