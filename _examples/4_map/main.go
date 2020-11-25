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
	printer := tableprinter.New(os.Stdout)
	// one header, many string values.
	books := map[string][]string{
		"Title": {
			"To Kill a Mockingbird (To Kill a Mockingbird) ",
			"The Hunger Games (The Hunger Games) ",
			"Harry Potter and the Order of the Phoenix (Harry Potter) ",
			"Pride and Prejudice ",
			"Animal Farm",
		},
	}

	/*
	  TITLE (5)
	 -----------------------------------------------------------
	  To Kill a Mockingbird (To Kill a Mockingbird)
	  The Hunger Games (The Hunger Games)
	  Harry Potter and the Order of the Phoenix (Harry Potter)
	  Pride and Prejudice
	  Animal Farm
	*/
	printer.Print(books)

	println()

	many := map[string][]person{
		"Access 1": {{"Georgios", "Callas"},
			{"Ioannis", "Christou"}},
		"Access 2": {
			{"Dimitrios", "Dellis"}},
		"Access 3": {{"Giannhs", "Christou"},
			{"Giwrgos", "Christou"},
			{"Oresths", "Christou"}},
		"Access 4": {{"Nikolaos", "Dellis"},
			{"Dionisis", "Dellis"}},
		"Access 5": {{"Fwths", "Papadopoulos"},
			{"Xrusostomos", "Papadopoulos"},
			{"Evriklia", "Papadopoulou"},
			{"Xrusa", "Papadopoulou"}},
	}

	/*
	  ACCESS 1           ACCESS 2           ACCESS 3           ACCESS 4          ACCESS 5
	 ------------------ ------------------ ------------------ ----------------- --------------------------
	  Georgios Callas    Dimitrios Dellis   Giannhs Christou   Nikolaos Dellis   Fwths Papadopoulos
	  Ioannis Christou                      Giwrgos Christou   Dionisis Dellis   Xrusostomos Papadopoulos
	                                        Oresths Christou                     Evriklia Papadopoulou
	                                                                             Xrusa Papadopoulou
	*/
	/*
	  ACCESS 3 (4)       ACCESS 4          ACCESS 5                   ACCESS 1           ACCESS 2
	 ------------------ ----------------- -------------------------- ------------------ ------------------
	  Giannhs Christou   Nikolaos Dellis   Fwths Papadopoulos         Georgios Callas    Dimitrios Dellis
	  Giwrgos Christou   Dionisis Dellis   Xrusostomos Papadopoulos   Ioannis Christou
	  Oresths Christou                     Evriklia Papadopoulou
	                                       Xrusa Papadopoulou
	*/
	printer.Print(many)

	println()

	onetoone := map[string]person{
		"Seller":   {"Georgios", "Callas"},
		"Consumer": {"Dimitrios", "Dellis"},
	}

	/*
	  SELLER            CONSUMER
	 ----------------- ------------------
	  Georgios Callas   Dimitrios Dellis
	*/
	printer.Print(onetoone)
}
