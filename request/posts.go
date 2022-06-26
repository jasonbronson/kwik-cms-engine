package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/repositories"

	"github.com/jasonbronson/kwik-cms-engine/request/response"
)

func GetPosts(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetPosts(db, params, "post")
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetPost(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("postid")
	var r *response.Response
	r = repositories.GetPostByID(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutPosts(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	var post model.Post
	if e := g.ShouldBindJSON(&post); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	e := repositories.UpdatePost(db, post)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func SetPosts(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	var post model.Post
	if e := g.ShouldBindJSON(&post); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	post.Type = "post"

	e := repositories.CreatePost(db, post)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func DeletePosts(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("postid")
	e := repositories.DeletePost(db, id)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func UpdatePublishDate(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("postid")

	var post model.Post
	if e := g.ShouldBindJSON(&post); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	e := repositories.UpdatePublishDate(db, post.PublishDate, id)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
