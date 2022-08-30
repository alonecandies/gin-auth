package repositories

import (
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository interface {
	InsertUser(user entities.User) entities.User
	UpdateUser(user entities.User) entities.User
	VerifyCredentials(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entities.User
	ProfileUser(userId uint64) entities.User
	MyBooks(userId uint64) []entities.Book
}

type userConnection struct {
	connection *gorm.DB
}

func NewUserRepository(connection *gorm.DB) UserRepository {
	return &userConnection{
		connection: connection,
	}
}

func hashPassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}

func (u *userConnection) InsertUser(user entities.User) entities.User {
	user.Password = hashPassword([]byte(user.Password))
	u.connection.Create(&user)
	return user
}

func (u *userConnection) UpdateUser(user entities.User) entities.User {
	if user.Password != "" {
		user.Password = hashPassword([]byte(user.Password))
	} else {
		var tempUser entities.User
		u.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}
	u.connection.Save(&user)
	return user
}

func (u *userConnection) VerifyCredentials(email string, password string) interface{} {
	var user entities.User
	res := u.connection.Where("email = ?", email).First(&user)
	if res.Error != nil {
		return nil
	} else {
		return user
	}
}

func (u *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	tx = u.connection.Where("email = ?", email).First(&entities.User{})
	return tx
}

func (u *userConnection) FindByEmail(email string) entities.User {
	var user entities.User
	u.connection.Where("email = ?", email).First(&user)
	return user
}

func (u *userConnection) ProfileUser(userId uint64) entities.User {
	var user entities.User
	u.connection.Where("id = ?", userId).First(&user)
	return user
}

func (u *userConnection) MyBooks(userId uint64) []entities.Book {
	var books []entities.Book
	u.connection.Where("user_id = ?", userId).Find(&books)
	return books
}