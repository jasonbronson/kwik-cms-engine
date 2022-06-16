package request

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/model"
	"github.com/jasonbronson/kwik-cms-engine/repositories"

	"github.com/jasonbronson/kwik-cms-engine/request/response"
	_ "github.com/jasonbronson/kwik-cms-engine/request/response"
)

func GetAuthors(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetAuthors(db, params)
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetAuthor(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("authorid")
	var r *response.Response
	r = repositories.GetAuthorByID(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutAuthors(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var author model.Author

	if e := g.ShouldBindJSON(&author); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}

	e := repositories.UpdateAuthor(db, author)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func SetAuthors(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var author model.Author

	if e := g.ShouldBindJSON(&author); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	e := repositories.CreateAuthor(db, author)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func DeleteAuthors(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	authorID := g.Param("authorid")
	e := repositories.DeleteAuthor(db, authorID)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
