package main

import (
	"os"

	"github.com/lensesio/tableprinter"
)

type (
	author struct {
		Name string `header:"name"`
		Age  int    `header:"age"`
		// Books -> header. count -> keyword to print the len(Books). none -> alternative value if len(Books) == 0.
		Books []book `header:"books,count,none"`

		// inline -> take the supposed `MyStruct`'s struct's fields tagged with headers and append those cells.
		// Field MyStruct `header:"inline"`
	}

	book struct {
		Title       string
		Description string
		Sales       int64
		Published   bool
	}
)

func main() {
	authors := []author{
		{"Author Name 1", 25, []book{{"Book Title 1", "Book Description 1", 132000, true}, {"Book Title 2", "Book Description 2", 164200, true}}},
		{"Author Name 2", 35, []book{{"Book Title 1 for Author 2", "Book Description 1 for Author 2", 0, false}}},
		{"Author Name 3", 42, nil},
		{"Author Name 4", 56, nil},
	}

	/*
	  NAME (4)        AGE   BOOKS
	 --------------- ----- -------
	  Author Name 1    25       2
	  Author Name 2    35       1
	  Author Name 3    42    none
	  Author Name 4    56    none
	*/
	tableprinter.Print(os.Stdout, authors) // prints to the "w" and returns the length of rows printed.
}
