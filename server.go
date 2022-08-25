package main

import (
	"github.com/alonecandies/mysql-gin-gorm-auth/api/configs/db"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/controllers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/repositories"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	conn           *gorm.DB                    = db.DBConnection()
	userRepository repositories.UserRepository = repositories.NewUserRepository(conn)
	jwtService     services.JWTService         = services.NewJWTService()
	authService    services.AuthService        = services.NewAuthService(userRepository)
	authController controllers.AuthController  = controllers.NewAuthController(authService, jwtService)
)

func main() {
	defer db.DBClose(conn)
	r := gin.Default()

	authRoutes := r.Group("api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	r.Run()
}
