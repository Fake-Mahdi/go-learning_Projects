package main

import "fmt"

type Book struct {
	author      string
	title       string
	isbn        string
	availabilty bool
}

var allBooks []Book

type User struct {
	name    string
	user_id string
	book    Book
}

var allUsers []User

func (book Book) add_book() Book {
	var author, title, isbn string
	var check int
	var availabilty bool

	fmt.Println("Give the Name of the author please:")
	fmt.Scanln(&author)

	fmt.Println("Give the title of the book please:")
	fmt.Scanln(&title)

	fmt.Println("Provide the ISBN please:")
	fmt.Scanln(&isbn)

	fmt.Printf("If available, provide 1; if not, provide 0: ")
	fmt.Scanln(&check)

	if check == 1 {
		availabilty = true
	} else {
		availabilty = false
	}

	new_book := Book{author: author, title: title, isbn: isbn, availabilty: availabilty}
	allBooks = append(allBooks, new_book)
	return new_book
}

func (book Book) borrow_book(title string, books []Book) {
	check := false
	for i, b := range books {
		if b.title == title && b.availabilty {
			fmt.Println("The book is available and now borrowed.")
			allBooks[i].availabilty = false
			check = true
			break
		}
		if b.title == title && !b.availabilty {
			fmt.Println("The book is in the library but not available now.")
			check = false
			break
		}
	}
	if !check {
		fmt.Println("The book is not available in the library.")
	}
}

func (book Book) is_available(title string, books []Book) bool {
	for _, b := range books {
		if b.title == title && b.availabilty {
			fmt.Println("The book is available.")
			return true
		}
		if b.title == title && !b.availabilty {
			fmt.Println("The book is in the library but not available now.")
			return false
		}
	}
	fmt.Println("The book is not available in the library.")
	return false
}

func (user *User) add_user(book Book) {
	var name, user_id string
	new_book := book

	fmt.Println("Add a client name here:")
	fmt.Scanln(&name)

	fmt.Println("Add a user ID:")
	fmt.Scanln(&user_id)

	new_user := User{name: name, user_id: user_id, book: new_book}
	allUsers = append(allUsers, new_user)
}

func (user *User) user_borrow_book(title string, b bookInterface) {
	if b.is_available(title, allBooks) {
		b.borrow_book(title, allBooks)
		fmt.Printf("%s borrowed the book: %s\n", user.name, title)
	} else {
		fmt.Printf("Sorry %s, the book %s is not available\n", user.name, title)
	}
}

func main() {
	var book Book
	var user User

	book1 := book.add_book()
	book2 := book.add_book()

	user.add_user(book1)
	user.add_user(book2)

	fmt.Println("\nAll books in library:")
	for _, b := range allBooks {
		fmt.Printf("Title: %s, Author: %s, Available: %t\n", b.title, b.author, b.availabilty)
	}

	fmt.Println("\nUser borrowing a book:")
	user.user_borrow_book(book1.title, book)

	fmt.Println("\nChecking availability after borrowing:")
	book.is_available(book1.title, allBooks)
}
