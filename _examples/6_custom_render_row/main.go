package main

import (
	"os"
	"reflect"

	"github.com/landoop/tableprinter"
)

// Person example.
type Person struct {
	Firstname string `header:"first name"`
	Lastname  string `header:"last name"`
}

func main() {
	printer := tableprinter.New(os.Stdout)

	person := Person{"Georgios", "Callas"}

	// Get headers manually from a struct value.
	v := reflect.ValueOf(person)
	headers := tableprinter.StructParser.ParseHeaders(v)

	// Render the headers, the rows and the positon of the number col.
	/*
	    FIRST NAME   LAST NAME
	   ------------ -----------
	*/
	printer.Render(headers, nil, nil, false)
	// The table is now rendered, but you can manually use the `RenderRow` to add more,
	// it will respect all the properties of the rendered table of this printer.

	// Get a single row based on that person.
	row, nums := tableprinter.StructParser.ParseRow(v)
	/*
	  Georgios     Callas
	*/
	printer.RenderRow(row, nums)

	// Add one more row for fun.
	row, nums = tableprinter.StructParser.ParseRow(reflect.ValueOf(Person{"Ioanis", "Christou"}))
	/*
	  Ioanis       Christou
	*/
	printer.RenderRow(row, nums)
}
