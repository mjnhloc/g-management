package middleware

import (
	"net/http"
	"strings"

	"g-management/pkg/dto"
	"g-management/pkg/shared/jwt"

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

		c.Set("user", claims)
		c.Next()
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

		roles := userClaims["https://api.gym-management.com/roles"].([]interface{})
		for _, role := range roles {
			if role == role {
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
