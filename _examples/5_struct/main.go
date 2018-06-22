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
	person := person{"Georgios", "Callas"}

	tableprinter.Print(os.Stdout, person)
}
