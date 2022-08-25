package middlewares

import (
	"net/http"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/helpers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			response := helpers.BuildErrorResponse(http.StatusUnauthorized, "Unauthorized", nil)
			c.AbortWithStatusJSON(response.Status, response)
			return
		}
		token, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			response := helpers.BuildErrorResponse(http.StatusUnauthorized, err.Error(), nil)
			c.AbortWithStatusJSON(response.Status, response)
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		response := helpers.BuildResponse(http.StatusOK, "Success", claims)
		c.JSON(response.Status, response)
		c.Next()
	}
}
