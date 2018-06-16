package tableprinter

import (
	"fmt"
	"io"
	"reflect"

	"github.com/olekukonko/tablewriter"
)

// Alignment is the alignment type (int).
//
// See `Printer#DefaultColumnAlignment` and `Printer#DefaultColumnAlignment` too.
type Alignment int

const (
	// AlignDefault is the default alignment (0).
	AlignDefault Alignment = iota
	// AlignCenter is the center aligment (1).
	AlignCenter
	// AlignRight is the right aligment (2).
	AlignRight
	// AlignLeft is the left aligment (3).
	AlignLeft
)

// Printer contains some information about the final table presentation.
// Look its `Print` function for more.
type Printer struct {
	AutoFormatHeaders bool
	AutoWrapText      bool

	BorderTop, BorderLeft, BorderRight, BorderBottom bool

	HeaderLine      bool
	HeaderAlignment Alignment

	RowLine         bool
	ColumnSeparator string
	NewLine         string
	CenterSeparator string

	DefaultAlignment Alignment // see `NumbersAlignment` too.
	NumbersAlignment Alignment

	RowLengthTitle func(int) bool
}

// Default is the default Table Printer.
var Default = Printer{
	AutoFormatHeaders: true,
	AutoWrapText:      true,

	BorderTop:    false,
	BorderLeft:   false,
	BorderRight:  false,
	BorderBottom: false,

	HeaderLine:      true,
	HeaderAlignment: AlignLeft,

	RowLine:         false, /* it could be true as well */
	ColumnSeparator: " ",
	NewLine:         "\n",
	CenterSeparator: " ", /* it could be empty as well */

	DefaultAlignment: AlignLeft,
	NumbersAlignment: AlignRight,

	RowLengthTitle: func(rowsLength int) bool {
		// if more than 3 then show the length of rows.
		return rowsLength > 3
	},
}

// Print calls and returns the result of the `Default.Print`,
// take a look at the `Printer#Print` function for details.
func Print(w io.Writer, v interface{}, filters ...interface{}) int {
	return Default.Print(w, v, filters...)
}

// Print usage:
// Print(writer, tt, func(t MyStruct) bool { /* or any type, depends on the type(s) of the "tt" */
// 	return t.Visibility != "hidden"
// })
//
// Returns the number of rows finally printed.
func (p *Printer) Print(w io.Writer, v interface{}, filters ...interface{}) int {
	table := tablewriter.NewWriter(w)
	table.SetAlignment(int(p.DefaultAlignment))

	var (
		headers             []string
		rows                [][]string
		numbersColsPosition []int
	)

	if val := reflect.Indirect(reflect.ValueOf(v)); val.Kind() == reflect.Slice {
		var f []RowFilter
		for i, n := 0, val.Len(); i < n; i++ {
			v := val.Index(i)

			if i == 0 {
				// make filters once instead of each time for each entry, they all have the same v type.
				f = MakeFilters(v, filters)
				headers = GetHeaders(v.Type())
			}

			if !v.IsValid() {
				rows = append(rows, []string{""})
				continue
			}

			right, row := GetRow(v)
			if i == 0 {
				numbersColsPosition = right
			}

			if CanAcceptRow(v, f) {
				rows = append(rows, row)
			}
		}
	} else {
		// single.
		headers = GetHeaders(val.Type())
		right, row := GetRow(val)
		numbersColsPosition = right
		if CanAcceptRow(val, MakeFilters(val, filters)) {
			rows = append(rows, row)
		}

	}

	if len(headers) == 0 {
		return 0
	}

	if p.RowLengthTitle != nil && p.RowLengthTitle(len(rows)) {
		headers[0] = fmt.Sprintf("%s (%d) ", headers[0], len(rows))
	}

	table.SetHeader(headers)
	table.AppendBulk(rows)

	table.SetAutoFormatHeaders(p.AutoFormatHeaders)
	table.SetAutoWrapText(p.AutoWrapText)
	table.SetBorders(tablewriter.Border{Top: p.BorderTop, Left: p.BorderLeft, Right: p.BorderRight, Bottom: p.BorderBottom})
	table.SetHeaderLine(p.HeaderLine)
	table.SetHeaderAlignment(int(p.HeaderAlignment))
	table.SetRowLine(p.RowLine)
	table.SetColumnSeparator(p.ColumnSeparator)
	table.SetNewLine(p.NewLine)
	table.SetCenterSeparator(p.CenterSeparator)

	columnAlignment := make([]int, len(headers), len(headers))
	for i := range columnAlignment {
		columnAlignment[i] = int(p.DefaultAlignment)

		for _, j := range numbersColsPosition {
			if i == j {
				columnAlignment[i] = int(p.NumbersAlignment)
				break
			}
		}

	}
	table.SetColumnAlignment(columnAlignment)

	table.Render()
	return len(rows)
}
