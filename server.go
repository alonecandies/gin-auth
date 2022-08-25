package main

import (
	"github.com/alonecandies/mysql-gin-gorm-auth/api/configs/db"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	conn           *gorm.DB                   = db.DBConnection()
	authController controllers.AuthController = controllers.NewAuthController()
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
