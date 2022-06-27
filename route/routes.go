package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/middlewares"
	"github.com/jasonbronson/kwik-cms-engine/request"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
	requestid "github.com/sumit-tembe/gin-requestid"
)

// Router func
func Router(newRelicApp *newrelic.Application) http.Handler {

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}
	router := gin.Default()
	router.RedirectTrailingSlash = true

	router.Use(nrgin.Middleware(newRelicApp))
	router.Use(cors.New(corsConfig))
	router.Use(requestid.RequestID(nil))

	router.GET("/", request.HealthCheck)
	router.GET("/healthz", request.HealthCheck)

	//Performance verify key on load forge
	loaderVerification := os.Getenv("LOAD_FORGE")
	if loaderVerification != "" {
		router.GET(fmt.Sprintf("/%v", loaderVerification), func(g *gin.Context) {
			g.Writer.WriteHeader(http.StatusOK)
			g.Writer.Write([]byte(loaderVerification))
		})
	}

	api := router.Group("/v1/api")
	{
		api.POST("/sign-in", request.PostLoginHandlerFunc)
		api.Use(middlewares.AuthMiddleware())
		dynamic := api.Group("/dynamicgroups")
		{
			dynamic.GET("", request.GetDynamicGroups)
			dynamic.GET("/:id", request.GetDynamicGroup)
			dynamic.POST("", request.PostDynamicGroup)
			dynamic.PUT("/:id", request.PutDynamicGroup)
			dynamic.DELETE("/:id", request.DeleteDynamicGroup)
		}
		users := api.Group("/users")
		{
			users.GET("", request.GetUsers)
			users.GET("/:userid", request.GetUser)
			users.POST("", request.PostUser)
			users.PUT("/:userid", request.PutUsers)
			users.DELETE("/:userid", request.DeleteUsers)
		}
		tags := api.Group("/tags")
		{
			tags.GET("", request.GetTags)
			tags.GET("/:tagid", request.GetTag)
			tags.POST("", request.SetTags)
			tags.PUT("/:tagid", request.PutTags)
			tags.DELETE("/:tagid", request.DeleteTags)
		}
		categories := api.Group("/categories")
		{
			categories.GET("", request.GetCategories)
			categories.GET("/:categoryid", request.GetCategory)
			categories.POST("", request.SetCategories)
			categories.PUT("/:categoryid", request.PutCategories)
			categories.DELETE("/:categoryid", request.DeleteCategories)
		}
		posts := api.Group("/posts")
		{
			posts.GET("", request.GetPosts)
			posts.GET("/:postid", request.GetPost)
			posts.POST("", request.SetPosts)
			posts.PUT("/:postid", request.PutPosts)
			posts.DELETE("/:postid", request.DeletePosts)
			posts.PUT("/publish/:id", request.UpdatePublishDate)
		}
		pages := api.Group("/pages")
		{
			pages.GET("", request.GetPages)
			pages.GET("/:pageid", request.GetPage)
			pages.POST("", request.SetPages)
			pages.PUT("/publish/:id", request.UpdatePublishDate)
			pages.PUT("/:pageid", request.PutPages)
			pages.DELETE("/:pageid", request.DeletePages)
		}
		media := api.Group("/media")
		{
			media.GET("", request.GetMedia)
			media.POST("", request.SetMedia)
			media.DELETE("/:mediaid", request.DeleteMedia)
		}
		roles := api.Group("/roles")
		{
			roles.GET("", request.GetRoles)
			roles.GET("/:roleid", request.GetRole)
			roles.POST("", request.SetRoles)
			roles.PUT("/:roleid", request.PutRoles)
			roles.DELETE("/:roleid", request.DeleteRoles)
		}
	}

	return router
}
