package controllers

import (
	"net/http"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/helpers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type authController struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func NewAuthController(authService services.AuthService, jwtService services.JWTService) AuthController {
	return &authController{
		authService: authService,
		jwtService:  jwtService,
	}
}

func (c *authController) Login(ctx *gin.Context) {
	var loginDTO dtos.LoginDTO
	errDTO := ctx.ShouldBind(&loginDTO)
	if errDTO != nil {
		response := helpers.BuildErrorResponse(http.StatusBadRequest, errDTO.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authRes := c.authService.VerifyCredentials(loginDTO.Email, loginDTO.Password)
	if v, ok := authRes.(entities.User); ok {
		token, err := c.jwtService.GenerateToken(uint(v.ID))
		if err != nil {
			response := helpers.BuildErrorResponse(http.StatusInternalServerError, err.Error(), helpers.EmptyResponse{})
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		response := helpers.BuildResponse(http.StatusOK, "Login successfully", token)
		ctx.JSON(http.StatusOK, response)
	} else {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, "Invalid credentials", helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}
}

func (c *authController) Register(ctx *gin.Context) {
	var registerDTO dtos.RegisterDTO
	errDTO := ctx.ShouldBind(&registerDTO)
	if errDTO != nil {
		response := helpers.BuildErrorResponse(http.StatusBadRequest, errDTO.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	if !c.authService.IsDuplicateEmail(registerDTO.Email) {
		response := helpers.BuildErrorResponse(http.StatusConflict, "Email already exists", helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusConflict, response)
		return
	} else {
		user := c.authService.CreateUser(&registerDTO)
		token, err := c.jwtService.GenerateToken(uint(user.ID))
		if err != nil {
			response := helpers.BuildErrorResponse(http.StatusInternalServerError, err.Error(), helpers.EmptyResponse{})
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		response := helpers.BuildResponse(http.StatusOK, "Register successfully", token)
		ctx.JSON(http.StatusOK, response)
	}
}
