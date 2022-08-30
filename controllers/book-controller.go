package controllers

import (
	"net/http"
	"strconv"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/helpers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/services"
	"github.com/gin-gonic/gin"
)

type BookController interface {
	InsertBook(c *gin.Context)
	UpdateBook(c *gin.Context)
	DeleteBook(c *gin.Context)
	AllBooks(c *gin.Context)
	FindBookById(c *gin.Context)
}

type bookController struct {
	bookService services.BookService
	jwtService services.JWTService
}

func NewBookController(bookService services.BookService, jwtService services.JWTService) BookController {
	return &bookController{
		bookService: bookService,
		jwtService: jwtService,
	}
}

func (bc *bookController) getUserIdByToken(token string) uint64 {
	aToken, err:= bc.jwtService.ValidateToken(token)
	if err != nil {
		panic(err)
	}
	claims:= aToken.Claims.(*services.JwtCustomClaims)
	return uint64(claims.UserId)
}

func (bc *bookController) InsertBook(c *gin.Context) {
	var bookCreateDTO dtos.BookCreateDTO
	errDTO:= c.ShouldBind(&bookCreateDTO)
	if errDTO != nil {
		res:= helpers.BuildErrorResponse(http.StatusBadRequest, "Invalid json", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	} else {
		token:= c.GetHeader("Authorization")
		userId:= bc.getUserIdByToken(token)
		bookCreateDTO.UserID = userId
		result:= bc.bookService.InsertBook(bookCreateDTO)
		response:= helpers.BuildResponse(http.StatusOK, "OK", result)
		c.JSON(http.StatusOK, response)
	}
}

func (bc *bookController) UpdateBook(c *gin.Context) {
	var bookUpdateDTO dtos.BookUpdateDTO
	errDTO:= c.ShouldBind(&bookUpdateDTO)
	if errDTO != nil {
		res:= helpers.BuildErrorResponse(http.StatusBadRequest, "Invalid json", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	} else {
		authHeader:= c.GetHeader("Authorization")
		token,err:= bc.jwtService.ValidateToken(authHeader)
		if err != nil {
			res:= helpers.BuildErrorResponse(http.StatusUnauthorized, "Unauthorized", helpers.EmptyResponse{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		claims:= token.Claims.(*services.JwtCustomClaims)
		bookUpdateDTO.UserID = uint64(claims.UserId)
		if (bc.bookService.IsAllowedToEdit(bookUpdateDTO.ID, bookUpdateDTO.UserID)) {
			result:= bc.bookService.UpdateBook(bookUpdateDTO)
			response:= helpers.BuildResponse(http.StatusOK, "OK", result)
			c.JSON(http.StatusOK, response)
		} else {
			res:= helpers.BuildErrorResponse(http.StatusForbidden, "Forbidden", helpers.EmptyResponse{})
			c.AbortWithStatusJSON(http.StatusForbidden, res)
		}
	}
}

func (bc *bookController) DeleteBook(c *gin.Context) {
	id,err:= strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		res:= helpers.BuildErrorResponse(http.StatusBadRequest, "Invalid id", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	authHeader:= c.GetHeader("Authorization")
	token,err:= bc.jwtService.ValidateToken(authHeader)
	if err != nil {
		res:= helpers.BuildErrorResponse(http.StatusUnauthorized, "Unauthorized", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusUnauthorized, res)
		return
	}
	claims:= token.Claims.(*services.JwtCustomClaims)
	userId := uint64(claims.UserId)
	if (bc.bookService.IsAllowedToEdit(id, userId)) {
		isDeleted:= bc.bookService.DeleteBook(id)
		if isDeleted {
			res:= helpers.BuildResponse(http.StatusOK, "Book deleted", helpers.EmptyResponse{})
			c.JSON(http.StatusOK, res)
		} else {
			res:= helpers.BuildErrorResponse(http.StatusBadRequest, "Failed to delete book", helpers.EmptyResponse{})
			c.AbortWithStatusJSON(http.StatusBadRequest, res)
		}
	} else {
		res:= helpers.BuildErrorResponse(http.StatusForbidden, "Forbidden", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusForbidden, res)
	}
}

func (bc *bookController) AllBooks(c *gin.Context) {
	var books []entities.Book = bc.bookService.AllBooks()
	res:= helpers.BuildResponse(http.StatusOK, "All books", books)
	c.JSON(http.StatusOK, res)
}

func (bc *bookController) FindBookById(c *gin.Context) {
	id,err:= strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		res:= helpers.BuildErrorResponse(http.StatusBadRequest, "Invalid id", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	var book entities.Book = bc.bookService.FindBookById(id)
	if (book == (entities.Book{})) {
		res:= helpers.BuildErrorResponse(http.StatusNotFound, "Book not found", helpers.EmptyResponse{})
		c.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	}
	res:= helpers.BuildResponse(http.StatusOK, "Book found", book)
	c.JSON(http.StatusOK, res)
}