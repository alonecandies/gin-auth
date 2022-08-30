package services

import (
	"log"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/repositories"
	"github.com/mashingan/smapping"
)

type UserService interface {
	Update(user dtos.UserUpdateDTO) entities.User
	Profile(userId uint64) entities.User
	MyBooks(userId uint64) []entities.Book
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (service *userService) Update(user dtos.UserUpdateDTO) entities.User {
	userToUpdate := entities.User{}
	err := smapping.FillStruct(&userToUpdate, smapping.MapFields(&user))
	if err != nil {
		log.Fatalf("Failed to map: %v", err)
	}
	updatedUser := service.userRepository.UpdateUser(userToUpdate)
	return updatedUser
}

func (service *userService) Profile(userId uint64) entities.User {
	return service.userRepository.ProfileUser(userId)
}

func (service *userService) MyBooks(userId uint64) []entities.Book {
	return service.userRepository.MyBooks(userId)
}
