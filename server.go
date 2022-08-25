package main

import (
  "github.com/gin-gonic/gin"
	"github.com/alonecandies/mysql-gin-gorm-auth/configs/db"
)

var (
	conn *gorm.DB = db.DBConnection()
	authController controllers.AuthController = controllers.NewAuthController()
)

func main() {
	defer conn.DBClose()
  r := gin.Default()
  
  authRoutes := r.Group("/auth")
  {
	 authRoutes.POST("/login")
  }

  r.Run()
}