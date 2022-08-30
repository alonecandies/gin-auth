package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/helpers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
)

type UserController interface {
	Update(c *gin.Context)
	Profile(c *gin.Context)
	MyBooks(c *gin.Context)
}

type userController struct {
	userService services.UserService
	jwtService  services.JWTService
}

func NewUserController(userService services.UserService, jwtService services.JWTService) UserController {
	return &userController{
		userService: userService,
		jwtService:  jwtService,
	}
}

func (c *userController) Update(ctx *gin.Context) {
	var user dtos.UserUpdateDTO
	errDTO := ctx.ShouldBind(&user)
	if errDTO != nil {
		response := helpers.BuildErrorResponse(http.StatusBadRequest, errDTO.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, errToken.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	claims := token.Claims.(*services.JwtCustomClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims.UserId), 10, 64)
	if err != nil {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, err.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	user.ID = id
	u := c.userService.Update(user)
	res := helpers.BuildResponse(http.StatusOK, "User updated", u)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) Profile(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, errToken.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	claims := token.Claims.(*services.JwtCustomClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims.UserId), 10, 64)
	if err != nil {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, err.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	u := c.userService.Profile(id)
	res := helpers.BuildResponse(http.StatusOK, "User profile", u)
	ctx.JSON(http.StatusOK, res)
}

func (c *userController) MyBooks(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, errToken.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	claims := token.Claims.(*services.JwtCustomClaims)
	id, err := strconv.ParseUint(fmt.Sprintf("%v", claims.UserId), 10, 64)
	if err != nil {
		response := helpers.BuildErrorResponse(http.StatusUnauthorized, err.Error(), helpers.EmptyResponse{})
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
		return
	}
	books := c.userService.MyBooks(id)
	res := helpers.BuildResponse(http.StatusOK, "User books", books)
	ctx.JSON(http.StatusOK, res)
}
