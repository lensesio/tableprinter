package main

import (
	"os"

	"github.com/kataras/tableprinter"
)

type person struct {
	FirstName string
	LastName  string
}

func (p person) String() string {
	return p.FirstName + " " + p.LastName
}

func main() {
	// one header, many string values.
	books := map[string][]string{
		"Title": []string{
			"To Kill a Mockingbird (To Kill a Mockingbird) ",
			"The Hunger Games (The Hunger Games) ",
			"Harry Potter and the Order of the Phoenix (Harry Potter) ",
			"Pride and Prejudice ",
			"Animal Farm",
		},
	}

	tableprinter.PrintMap(os.Stdout, books)

	println()

	many := map[string][]person{
		"Access 1": []person{{"Georgios", "Callas"},
			{"Ioannis", "Christou"}},
		"Access 2": []person{
			{"Dimitrios", "Dellis"}},
		"Access 3": []person{{"Dimitrios3", "Dellis3"},
			{"Nikolaos3", "Doukas3"},
			{"Third3", "Name3"}},
		"Access 4": []person{{"Nikolaos", "Doukas"},
			{"Third", "Name"}},
	}

	/*
	  ACCESS 1           ACCESS 2           ACCESS 3             ACCESS 4
	 ------------------ ------------------ -------------------- -----------------
	  Georgios Callas    Dimitrios Dellis   Dimitrios3 Dellis3   Nikolaos Doukas
	  Ioannis Christou                      Nikolaos3 Doukas3    Third Name
	                                        Third3 Name3
	*/
	tableprinter.PrintMap(os.Stdout, many)

	println()

	onetoone := map[string]person{
		"Seller":   person{"Georgios", "Callas"},
		"Consumer": person{"Dimitrios", "Dellis"},
	}

	tableprinter.PrintMap(os.Stdout, onetoone)
}
