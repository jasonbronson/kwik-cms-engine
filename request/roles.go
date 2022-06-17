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

func GetRoles(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetRoles(db, params)
	response.Standard(g.Writer, http.StatusOK, r)
}

func GetRole(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("roleid")
	r := repositories.GetRole(db, id)
	response.Standard(g.Writer, http.StatusOK, r)
}

func PutRoles(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var role model.Role
	if e := g.ShouldBindJSON(&role); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	e := repositories.UpdateRole(db, role)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}

func SetRoles(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	var role model.Role
	if e := g.ShouldBindJSON(&role); e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	id, e := repositories.CreateRole(db, role)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, fmt.Sprint(id))
}

func DeleteRoles(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	roleID := g.Param("roleid")
	e := repositories.DeleteRole(db, roleID)
	if e != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, e)
		return
	}
	response.Action(g.Writer, http.StatusOK, "success")
}
