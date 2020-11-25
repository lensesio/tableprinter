package main

import (
	"os"

	"github.com/lensesio/tableprinter"
)

type person struct {
	Firstname string `header:"first name"`
	Lastname  string `header:"last name"`
}

func main() {
	person := person{"Georgios", "Callas"}

	/*
	  FIRST NAME   LAST NAME
	 ------------ -----------
	  Georgios     Callas
	*/
	tableprinter.Print(os.Stdout, person)
}
