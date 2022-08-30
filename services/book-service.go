package services

import (
	"log"

	"github.com/alonecandies/mysql-gin-gorm-auth/api/dtos"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/entities"
	"github.com/alonecandies/mysql-gin-gorm-auth/api/repositories"
	"github.com/mashingan/smapping"
)

type BookService interface {
	InsertBook(book dtos.BookCreateDTO) entities.Book
	UpdateBook(book dtos.BookUpdateDTO) entities.Book
	DeleteBook(id uint64) bool
	AllBooks() []entities.Book
	FindBookById(id uint64) entities.Book
	IsAllowedToEdit(userID uint64, bookID uint64) bool
}

type bookService struct {
	bookRepository repositories.BookRepository
}

func NewBookService(bookRepository repositories.BookRepository) BookService {
	return &bookService{
		bookRepository: bookRepository,
	}
}

func (bs *bookService) InsertBook(b dtos.BookCreateDTO) entities.Book {
	book := entities.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatal(err)
	}
	res := bs.bookRepository.InsertBook(book)
	return res
}

func (bs *bookService) UpdateBook(b dtos.BookUpdateDTO) entities.Book {
	book := entities.Book{}
	err := smapping.FillStruct(&book, smapping.MapFields(&b))
	if err != nil {
		log.Fatal(err)
	}
	res := bs.bookRepository.UpdateBook(book)
	return res
}

func (bs *bookService) DeleteBook(id uint64) bool {
	return bs.bookRepository.DeleteBook(id)
}

func (bs *bookService) AllBooks() []entities.Book {
	return bs.bookRepository.AllBooks()
}

func (bs *bookService) FindBookById(id uint64) entities.Book {
	return bs.bookRepository.FindBookById(id)
}

func (bs *bookService) IsAllowedToEdit(userID uint64, bookID uint64) bool {
	return bs.bookRepository.FindBookById(bookID).User.ID == userID
}
