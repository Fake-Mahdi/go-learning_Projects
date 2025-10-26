package main

type bookInterface interface {
	add_book() Book
	borrow_book(title string, books []Book)
	is_available(title string, books []Book) bool
}
