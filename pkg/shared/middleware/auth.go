package middleware

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"g-management/pkg/dto"
	"g-management/pkg/shared/jwt"
	"g-management/pkg/shared/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	jwtGo "github.com/golang-jwt/jwt/v5"
)

func CheckAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		dataErrorUnauthorized := &dto.BaseErrorResponse{
			Error: &dto.ErrorResponse{
				Message: map[string]interface{}{
					"access_token": "Error accessing the resource, please provide a valid access token",
				},
			},
		}

		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, dataErrorUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, jwt.BearerAuthorizationPrefix)
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, dataErrorUnauthorized)
			return
		}

		decodedToken, err := jwt.Decode(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dataErrorUnauthorized)
			return
		}

		claims, ok := decodedToken.Claims.(jwtGo.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dataErrorUnauthorized)
			return
		}

		sub, subOK := claims["sub"].(string)
		iss, issOK := claims["iss"].(string)
		if !subOK || !issOK {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dataErrorUnauthorized)
			return
		}

		auth0Domain := os.Getenv("AUTH0_DOMAIN")
		audience := os.Getenv("AUTH0_AUDIENCE")
		if iss != auth0Domain {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dataErrorUnauthorized)
			return
		}
		c.Set("auth0_user_id", sub)

		middleware := getMiddlewareAuth0(c, auth0Domain, audience)
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			c.Request = r
		}

		middleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)
		if encounteredError {
			return
		}
		// token := r.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		// mapClaims := token.CustomClaims.(*validator.MapClaims)

		// c.Set("user", *mapClaims) // set toàn bộ claims để dùng sau

		c.Next()
		return
	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": "User claims not found, please provide a valid access token",
					},
				},
			})
			return
		}

		userClaims := claims.(jwtGo.MapClaims)
		userRoleUrl := os.Getenv("AUTH0_AUDIENCE") + "/roles"
		roles, ok := userClaims[userRoleUrl].([]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: map[string]interface{}{
						"access_token": "User roles not found, please provide a valid access token",
					},
				},
			})
			return
		}

		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, &dto.BaseErrorResponse{
			Error: &dto.ErrorResponse{
				Message: map[string]interface{}{
					"access_token": "You do not have permission to access this resource",
				},
			},
		})
	}
}

func getMiddlewareAuth0(c *gin.Context, domain, audience string) *jwtmiddleware.JWTMiddleware {
	issuerURL, err := url.Parse("https://" + domain + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{audience},
		validator.WithAllowedClockSkew(30*time.Second),
	)
	if err != nil {
		log.Fatalf("failed to set up the validator: %v", err)
	}

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				c.AbortWithStatusJSON(http.StatusRequestTimeout, &dto.BaseErrorResponse{
					Error: &dto.ErrorResponse{
						Message: utils.ErrorTokenExpiredMsgJp,
					},
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.BaseErrorResponse{
				Error: &dto.ErrorResponse{
					Message: err.Error(),
				},
			})
		}
	}

	return jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)
}
