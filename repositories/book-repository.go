package repositories

import (
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"gorm.io/gorm"
)

type BookRepository interface {
	InsertBook(book entities.Book) entities.Book
	UpdateBook(book entities.Book) entities.Book
	DeleteBook(id uint64) bool
	AllBooks() []entities.Book
	FindBookById(id uint64) entities.Book
}

type bookConnection struct {
	connection *gorm.DB
}

func NewBookRepository(connection *gorm.DB) BookRepository {
	return &bookConnection{
		connection: connection,
	}
}

func (db *bookConnection) InsertBook(book entities.Book) entities.Book {
	db.connection.Create(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) UpdateBook(book entities.Book) entities.Book {
	db.connection.Save(&book)
	db.connection.Preload("User").Find(&book)
	return book
}

func (db *bookConnection) DeleteBook(id uint64) bool {
	book := entities.Book{}
	db.connection.First(&book, id)
	if book.ID == 0 {
		return false
	}
	db.connection.Delete(&book)
	return true
}

func (db *bookConnection) AllBooks() []entities.Book {
	var books []entities.Book
	db.connection.Preload("User").Find(&books)
	return books
}

func (db *bookConnection) FindBookById(id uint64) entities.Book {
	var book entities.Book
	db.connection.First(&book, id)
	return book
}