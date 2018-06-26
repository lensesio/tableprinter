package main

import (
	"os"

	"github.com/kataras/tableprinter"
)

type person struct {
	Firstname string `header:"first name"`
	Lastname  string `header:"last name"`
}

func main() {
	persons := []person{
		{"Chris", "Doukas"},
		{"Georgios", "Callas"},
		{"Ioannis", "Christou"},
		{"Dimitrios", "Dellis"},
		{"Nikolaos", "Doukas"},
	}

	onlyDoukasFilter := func(p person) bool {
		return p.Lastname == "Doukas"
	}

	/*
	  FIRST NAME   LAST NAME
	 ------------ -----------
	  Chris        Doukas
	  Nikolaos     Doukas
	*/
	tableprinter.Print(os.Stdout, persons, onlyDoukasFilter)
}
