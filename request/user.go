package request

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/jasonbronson/kwik-cms-engine/model"

	// uuid "github.com/satori/go.uuid"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	"github.com/jasonbronson/kwik-cms-engine/repositories"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	_ "github.com/jasonbronson/kwik-cms-engine/request/response"
	"golang.org/x/crypto/bcrypt"
)

func GetUsers(g *gin.Context) {

	db := config.Cfg.GormDB
	db = db.WithContext(g)
	params := defaultParameters(g)
	r := repositories.GetUsers(db, params)
	g.JSON(http.StatusOK, r)
}

func GetUser(g *gin.Context) {
	db := config.Cfg.GormDB
	db = db.WithContext(g)
	id := g.Param("userid")
	r := repositories.GetUser(db, id)
	g.JSON(http.StatusOK, r)
}

func PostLoginHandlerFunc(g *gin.Context) {
	// instrument databases with context for new relic
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	jwtConfig := config.Cfg.JwtConfig

	if email, password, ok := g.Request.BasicAuth(); ok {
		if email != "" && password != "" {
			fmt.Println("user", email, password)
			user := repositories.GetUserByEmail(db, email)
			if user.ID == "" {
				response.ErrorResponse(g.Writer, http.StatusUnauthorized, errors.New("username and/or password was incorrect 1"))
				return
			}
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
			if err != nil {
				response.ErrorResponse(g.Writer, http.StatusUnauthorized, errors.New("username and/or password was incorrect 2"))
				return
			}

			idV4, err := uuid.NewV4()

			claims := helpers.CustomClaims{
				StandardClaims: jwt.StandardClaims{
					Audience: jwtConfig.Audience,
					Id:       idV4.String(),
					IssuedAt: time.Now().Unix(),
					Issuer:   jwtConfig.Issuer,
					Subject:  "ddd",
				},
				Scope:      "user",
				Email:      email,
				Expiration: nil,
				Username:   "",
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signedString, err := token.SignedString([]byte(jwtConfig.Secret))
			if err != nil {
				response.ErrorResponse(g.Writer, http.StatusUnauthorized, err)
				return
			}

			g.Writer.Header().Set("Content-Type", "text/plain")
			g.Writer.WriteHeader(http.StatusOK)
			g.Writer.Write([]byte(signedString))

			return
		} else {
			response.ErrorResponse(g.Writer, http.StatusUnauthorized, errors.New("username and/or password was empty"))
			return
		}
	} else {
		response.ErrorResponse(g.Writer, http.StatusUnauthorized, errors.New("username and password weren't supplied"))
		return
	}

}

func PostUser(g *gin.Context) {
	// instrument databases with context for new relic
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	// 1.Get Body parameters
	var user model.User
	err := g.BindJSON(&user)
	if err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}
	user.Password = string(hashedPassword)
	// 2. check email
	userFromDB := repositories.GetUserByEmail(db, user.Email)
	if len(userFromDB.Email) > 0 {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, fmt.Errorf("email already exists"))
		return
	}

	// 3.Create User
	if err := repositories.CreateUser(db, user); err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}

	response.Action(g.Writer, http.StatusOK, "success")

}

func PutUsers(g *gin.Context) {
	// instrument databases with context for new relic
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	// 1.Get Body parameters
	userID := g.Param("userid")

	var user model.User

	err := g.BindJSON(&user)
	if err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}
	if userID != user.ID {
		response.ErrorResponse(g.Writer, http.StatusBadRequest, fmt.Errorf("patch request missing resource id: %v", userID))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}
	user.Password = string(hashedPassword)
	// 2.Update User
	if err := repositories.UpdateUser(db, user); err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}

	response.Action(g.Writer, http.StatusOK, "SUCCESS")

}

func DeleteUsers(g *gin.Context) {

	// instrument databases with context for new relic
	db := config.Cfg.GormDB
	db = db.WithContext(g)

	// 1.Get Context Params
	userID := g.Param("userid")

	err := repositories.DeleteUser(db, userID)
	if err != nil {
		response.ErrorResponse(g.Writer, http.StatusInternalServerError, err)
		return
	}

	response.Action(g.Writer, http.StatusOK, "SUCCESS")

}
