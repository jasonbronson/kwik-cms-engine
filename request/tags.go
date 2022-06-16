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

func GetTags(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetTags(db, params)
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetTag(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("tagid")
	r := repositories.GetTag(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutTags(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var tag model.Tag
	if e := g.ShouldBindJSON(&tag); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	e := repositories.UpdateTag(db, tag)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func SetTags(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var tag model.Tag
	if e := g.ShouldBindJSON(&tag); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	id, e := repositories.CreateTag(db, tag)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, fmt.Sprint(id))
}

func DeleteTags(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	tagID := g.Param("tagid")
	e := repositories.DeleteTag(db, tagID)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
