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
		"Access 3": []person{{"Giannhs", "Christou"},
			{"Giwrgos", "Christou"},
			{"Oresths", "Christou"}},
		"Access 4": []person{{"Nikolaos", "Dellis"},
			{"Dionisis", "Dellis"}},
		"Access 5": []person{{"Fwths", "Papadopoulos"},
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
	tableprinter.PrintMap(os.Stdout, many)

	println()

	onetoone := map[string]person{
		"Seller":   person{"Georgios", "Callas"},
		"Consumer": person{"Dimitrios", "Dellis"},
	}

	tableprinter.PrintMap(os.Stdout, onetoone)
}
