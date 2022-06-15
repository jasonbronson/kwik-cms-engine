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
		users := api.Group("/users")
		{
			users.GET("", request.GetUsers)
			users.GET("/:userid", request.GetUser)
			users.POST("", request.PostUser)
			users.PUT("/:userid", request.PutUsers)
			users.DELETE("/:userid", request.DeleteUsers)
		}
	}

	return router
}
