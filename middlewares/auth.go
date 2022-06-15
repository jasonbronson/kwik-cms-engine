package middlewares

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	"github.com/jasonbronson/kwik-cms-engine/library/helpers"
	cxt "github.com/jasonbronson/kwik-cms-engine/model/context"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(g *gin.Context) {
		jwtConfig := config.Cfg.JwtConfig
		// 0. Setup NewRelic segment and instrument Redis / Gorm
		txn := newrelic.FromContext(g)
		segment := txn.StartSegment("v2 AuthMiddleware")
		tryEndSegment := func() {
			if segment != nil {
				segment.End()
				segment = nil
			}
		}
		defer tryEndSegment()
		endSegmentAndServeHTTP := func(w http.ResponseWriter, r *http.Request) {
			tryEndSegment()
			g.Next()
		}

		// 1. Parse token and get token text
		tokenText, err := helpers.GetTokenFromRequest(g.Request)
		if err != nil {
			g.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrResponse{StatusCode: http.StatusUnauthorized, Title: "cannot get token from request", Message: err.Error()})
			return
		}

		if tokenText == "" {
			log.Println("Bearer token is required")
			g.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrResponse{StatusCode: http.StatusUnauthorized, Title: "Bearer token is required"})
			return
		}

		// 2. Special treatment for anonymous account
		if helpers.IsAnonymousAccountType(tokenText) {
			if g.Request.Method != "GET" {
				log.Println("IsAnonymousAccountType attempting to use Non GET methods")
				response.ErrorResponse(g.Writer, http.StatusUnauthorized, errors.New("API not authorized for this method"))
				g.Abort()
				return
			}
			expiration := time.Now().Add(time.Duration(300) * time.Minute)
			cc := helpers.CustomClaims{
				Expiration: &expiration,
			}
			ctx := helpers.SetContext(cxt.ContextCustomClaims, cc, g)
			// g.Writer.Header().Set("authtoken", "")
			// g.Writer.Header().Set("Access-Control-Expose-Headers", "authtoken")
			endSegmentAndServeHTTP(g.Writer, g.Request.WithContext(ctx))
			return
		}

		// 3. Get the token
		token, _ := jwt.ParseWithClaims(tokenText, &helpers.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		})

		if token == nil {
			log.Println("Token is not parsable ", tokenText)
			response.ErrorResponse(g.Writer, http.StatusUnauthorized, errors.New("token is not parsable"))
			g.Abort()
			return
		}

		if claims, ok := token.Claims.(*helpers.CustomClaims); ok && token.Valid {

			//4. Verify integrity of token
			err = helpers.VerifyClaims(claims, jwtConfig)
			if err != nil {
				log.Println("AuthMiddleware: VerifyClaims token failed ", claims)
				g.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrResponse{StatusCode: http.StatusUnauthorized, Title: "invalid token", Message: err.Error()})
				return
			}

			// 9.set the auth context
			helpers.SetContext(cxt.ContextCustomClaims, *claims, g)
			endSegmentAndServeHTTP(g.Writer, g.Request)
			return
		} else {
			log.Println("Parsing token failed ", token.Claims.Valid())
			g.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrResponse{StatusCode: http.StatusUnauthorized, Title: "invalid token", Message: "parsing token failure"})
			return
		}

	}
}
