package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/repositories"

	"github.com/jasonbronson/kwik-cms-engine/request/response"
)

func GetPages(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetPosts(db, params, "page")
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetPage(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("pageid")
	var r *response.Response
	r = repositories.GetPostByID(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutPages(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	var page model.Post
	if e := g.ShouldBindJSON(&page); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	e := repositories.UpdatePost(db, page)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func SetPages(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	var page model.Post
	if e := g.ShouldBindJSON(&page); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	page.Type = "page"
	e := repositories.CreatePost(db, page)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func DeletePages(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("pageid")
	e := repositories.DeletePost(db, id)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
