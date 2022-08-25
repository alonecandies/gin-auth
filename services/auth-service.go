package services

import (
	"log"
	"net/http"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/helpers"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/repositories"
	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredentials(email string, password string) interface{}
	CreateUser(userCreateDTO *dtos.UserCreateDTO) entities.User
	FindByEmail(email string) entities.User
	IsDuplicateEmail(email string) bool
}

type authService struct {
	userRepository repositories.UserRepository
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

func comparePassword(password []byte, hash []byte) bool {
	return bcrypt.CompareHashAndPassword(hash, password) == nil
}

func (service *authService) VerifyCredentials(email string, password string) interface{} {
	res := service.userRepository.VerifyCredentials(email, password)
	if v, ok := res.(entities.User); ok {
		comparePassword := comparePassword([]byte(password), []byte(v.Password))
		if v.Email == email && comparePassword {
			return helpers.BuildResponse(http.StatusOK, "Success", v)
		} else {
			return helpers.BuildErrorResponse(http.StatusUnauthorized, "Unauthorized", nil)
		}
	} else {
		return res
	}
}

func (service *authService) CreateUser(userCreateDTO *dtos.UserCreateDTO) entities.User {
	user := entities.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	res := service.userRepository.InsertUser(user)
	return res
}

func (service *authService) FindByEmail(email string) entities.User {
	return service.userRepository.FindByEmail(email)
}

func (service *authService) IsDuplicateEmail(email string) bool {
	res := service.userRepository.IsDuplicateEmail(email)
	return !(res.Error == nil)
}
