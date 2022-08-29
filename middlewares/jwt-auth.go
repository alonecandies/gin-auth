package middlewares

import (
	"log"
	"net/http"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/helpers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
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
		if token.Valid {
			claims := token.Claims.(*services.JwtCustomClaims)
			log.Println("Claim[user_id]: ", claims.UserId)
			log.Println("Claim[issuer] :", claims.Issuer)
		} else {
			log.Println(err)
			response := helpers.BuildErrorResponse(http.StatusBadRequest, err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
	}
}
