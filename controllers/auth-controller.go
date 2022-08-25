package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	
}

func NewAuthController() AuthController {
	return &authController{}
}

func (c *authController) Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Helllo login",
	})
}

func (c *authController) Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Helllo register",
	})
}
