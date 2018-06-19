package main

import (
	"os"

	"github.com/kataras/tableprinter"
)

type (
	firstThing struct {
		Name string `header:"Name"`
		Age  int    `header:"Age"`
	}

	secondThing struct {
		Title       string `header:"Title"`
		Description string `header:"Description"`
	}
)

func main() {
	things := []interface{}{firstThing{"First Thing Name", 25}, secondThing{"Second Thing Title", "Second Thing Description"}}

	// TODO:
	tableprinter.Print(os.Stdout, things)
}
