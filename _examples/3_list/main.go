package main

import (
	"os"

	"github.com/lensesio/tableprinter"
)

type person struct {
	FirstName string
	LastName  string
}

func (p person) String() string {
	return p.FirstName + " " + p.LastName
}

func main() {
	books := []string{
		"To Kill a Mockingbird (To Kill a Mockingbird) ",
		"The Hunger Games (The Hunger Games) ",
		"Harry Potter and the Order of the Phoenix (Harry Potter) ",
		"Pride and Prejudice ",
		"Animal Farm",
	}

	/*
	  BOOKS (5)
	 -----------------------------------------------------------
	  To Kill a Mockingbird (To Kill a Mockingbird)
	  The Hunger Games (The Hunger Games)
	  Harry Potter and the Order of the Phoenix (Harry Potter)
	  Pride and Prejudice
	  Animal Farm
	*/
	tableprinter.PrintHeadList(os.Stdout, books, "Books")

	println()

	numbers := []int{13213, 24554, 376575, 4321321321321, 5654654, 6654654, 787687, 8876876, 9321321}

	/*
	  NUMBERS (9)
	 --------------
	         13.2K
	         24.5K
	        376.5K
	          4.3T
	          5.6M
	          6.6M
	        787.6K
	          8.8M
	          9.3M
	*/
	tableprinter.PrintHeadList(os.Stdout, numbers, "Numbers")

	println()

	// DISCLAIMER: those are imaginary persons.
	persons := []person{
		{"Georgios", "Callas"},
		{"Ioannis", "Christou"},
		{"Dimitrios", "Dellis"},
		{"Nikolaos", "Doukas"},
	}

	/*
	  PERSONS (4)
	 ------------------
	  Georgios Callas
	  Ioannis Christou
	  Dimitrios Dellis
	  Nikolaos Doukas
	*/
	tableprinter.PrintHeadList(os.Stdout, persons, "Persons")
}
