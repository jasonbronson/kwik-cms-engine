package helpers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jasonbronson/kwik-cms-engine/config"
	cxt "github.com/jasonbronson/kwik-cms-engine/model/context"
	"github.com/jasonbronson/kwik-cms-engine/request/response"
)

const (
	AccessTokenMagicNumber = "8977"
	HttpStatusTokenExpired = 452
)

type CustomClaims struct {
	jwt.StandardClaims
	Scope      string     `json:"scope"`
	Email      string     `json:"email"`
	Expiration *time.Time `json:"expiration"`
	Username   string     `json:"user_name"`
}

func IsAnonymousAccountType(tokenText string) bool {
	return strings.HasPrefix(tokenText, AccessTokenMagicNumber)
}

func VerifyClaims(claims *CustomClaims, jwtConfig *config.JWTConfig) error {
	if !claims.VerifyIssuer(jwtConfig.Issuer, true) {
		return errors.New("invalid JWT Issuer claim")
	}
	if !claims.VerifyAudience(jwtConfig.Audience, true) {
		return errors.New("invalid JWT Audience claim")
	}
	return nil
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	tokenString := r.Header.Get("Authorization")
	splitToken := strings.Split(tokenString, "Bearer")
	if len(splitToken) != 2 {
		return "", errors.New("auth token incorrect or was't supplied")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

func SetContext(name cxt.ContextKey, claims CustomClaims, g *gin.Context) *gin.Context {
	g.Set(string(name), claims)
	return g
}

func GetClaimsFromContext(g *gin.Context) *CustomClaims {
	claimsCtx := g.Value(string(cxt.ContextCustomClaims)).(CustomClaims)
	return &claimsCtx
}

func GetClaimsFromRequest(g *gin.Context) *CustomClaims {
	tokenText, err := GetTokenFromRequest(g.Request)
	if err != nil {
		response.ErrorResponse(g.Writer, http.StatusUnauthorized, err)
		return nil
	}
	token, _ := jwt.ParseWithClaims(tokenText, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JwtConfig.Secret), nil
	})
	claims, _ := token.Claims.(*CustomClaims)
	return claims
}

func GetUserId(g *gin.Context) string {
	if g.Value(string(cxt.ContextCustomClaims)) == nil {
		return ""
	}
	claims := GetClaimsFromContext(g)
	return claims.Subject
}

func IsTokenExpired(tokenString string, jwtConfig *config.JWTConfig) (isExpired bool) {
	claim := GetCustomClaimFromString(tokenString, jwtConfig)
	if claim != nil && claim.Expiration != nil {
		if claim.Expiration.Before(time.Now()) {
			isExpired = true
		}
	}
	return
}

func GetUserIdFromToken(tokenString string, jwtConfig *config.JWTConfig) (id string) {
	claim := GetCustomClaimFromString(tokenString, jwtConfig)
	if claim != nil && claim.Subject != "" {
		return claim.Subject
	}
	return
}

func GetCustomClaimFromString(tokenString string, jwtConfig *config.JWTConfig) *CustomClaims {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtConfig.Secret), nil
	})
	if err != nil {
		return nil
	}
	return token.Claims.(*CustomClaims)
}

func GetSubscriberEntitlement(g *gin.Context) *CustomClaims {
	claims := GetClaimsFromContext(g)
	return claims
}

func GetEmail(g *gin.Context) string {
	if g.Value(string(cxt.ContextCustomClaims)) == nil {
		return ""
	}
	claims := GetClaimsFromContext(g)

	return claims.Email
}

// func NewToken(ctx context.Context, s *db_model.PropertySubscriber) string {
// 	return RenewToken(ctx, s, uuid.NewV4().String())
// }

// func RenewToken(g context.Context, s *db_model.PropertySubscriber, claimID string) string {
// 	jwtConfig := config.JWTConfig{}
// 	jwtConfig.Secret = config.Cfg.JwtConfig.Secret
// 	jwtConfig.Issuer = config.Cfg.JwtConfig.Issuer
// 	jwtConfig.Audience = config.Cfg.JwtConfig.Audience
// 	jwtExpirationInMins, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
// 	expiration := time.Now().Add(time.Minute * time.Duration(jwtExpirationInMins))
// 	username := ""
// 	if s.UserName != nil {
// 		username = *s.UserName
// 	}
// 	claims := CustomClaims{
// 		StandardClaims: jwt.StandardClaims{
// 			Audience: jwtConfig.Audience,
// 			Id:       claimID,
// 			IssuedAt: time.Now().Unix(),
// 			Issuer:   jwtConfig.Issuer,
// 			Subject:  s.Id,
// 		},
// 		Scope:      JWTSubscriber,
// 		Email:      s.Email,
// 		Expiration: &expiration,
// 		Username:   username,
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedString, err := token.SignedString([]byte(config.Cfg.JwtConfig.Secret))
// 	if err != nil {
// 		log.Println("Error signing token ")
// 		return ""
// 	}
// 	return signedString
// }
