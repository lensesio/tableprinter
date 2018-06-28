package main

import (
	"os"
	"sort"

	"github.com/kataras/tablewriter"
	"github.com/landoop/tableprinter"
)

type person struct {
	Firstname string `header:"first name"`
	Lastname  string `header:"last name"`
}

func main() {
	printer := tableprinter.New(os.Stdout)
	persons := []person{
		{"Chris", "Doukas"},
		{"Georgios", "Callas"},
		{"Ioannis", "Christou"},
		{"Nikolaos", "Doukas"},
		{"Dimitrios", "Dellis"},
	}

	sort.Slice(persons, func(i, j int) bool {
		return persons[j].Firstname > persons[i].Firstname
	})

	/*
		│─────────────────│───────────│
		│ FIRST NAME (5)  │ LAST NAME │ <- Green letters, black background header box.
		│─────────────────│───────────│
		│ Chris           │ Doukas    │
		│ Dimitrios       │ Dellis    │
		│ Georgios        │ Callas    │
		│ Ioannis         │ Christou  │
		│ Nikolaos        │ Doukas    │
		│─────────────────│───────────│
	*/

	// printer.HeaderLine = false // to disable headers.
	// printer.DefaultAlignment to change the alignment of cells.
	// printer.HeaderAlignment to change the alignment of header text.
	// printer.HeaderColors to set colors for each header manually, must match the number of headers.
	// printer.RowLengthTitle = func(n int) bool {
	// 	return n > 4
	// } // to change if and when the number of total rows should be shown after the first header, defaults to > 3.
	// printer.AutoWrapText = true // to enable cell's text wrap.
	// printer.NewLine = "\n" // to modify the new line for cells.
	printer.BorderTop, printer.BorderBottom, printer.BorderLeft, printer.BorderRight = true, true, true, true
	printer.CenterSeparator = "│"
	printer.ColumnSeparator = "│"
	printer.RowSeparator = "─"
	printer.HeaderBgColor = tablewriter.BgBlackColor // set header background color for all headers.
	printer.HeaderFgColor = tablewriter.FgGreenColor // set header foreground color for all headers.
	printer.Print(persons)
}
