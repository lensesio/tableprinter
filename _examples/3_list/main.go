package main

import (
	"os"

	"github.com/kataras/tableprinter"
)

func main() {
	books := []string{
		"To Kill a Mockingbird (To Kill a Mockingbird) ",
		"The Hunger Games (The Hunger Games) ",
		"Harry Potter and the Order of the Phoenix (Harry Potter) ",
		"Pride and Prejudice ",
		"Animal Farm",
	}

	_ = books
	/// TODO:
	// tableprinter.PrintHeadList(os.Stdout, books)
}
