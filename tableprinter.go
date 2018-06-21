package tableprinter

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"

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
	// out can not change during its work because the `acquire/release table` must work with only one output source,
	// a new printer should be declared for a different output.
	out io.Writer

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

	pool sync.Pool
}

// Default is the default Table Printer.
var Default = Printer{
	out:               os.Stdout,
	AutoFormatHeaders: true,
	AutoWrapText:      false,

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

func New(w io.Writer) *Printer {
	return &Printer{
		out: w,

		AutoFormatHeaders: Default.AutoFormatHeaders,
		AutoWrapText:      Default.AutoWrapText,

		BorderTop:    Default.BorderTop,
		BorderLeft:   Default.BorderLeft,
		BorderRight:  Default.BorderRight,
		BorderBottom: Default.BorderBottom,

		HeaderLine:      Default.HeaderLine,
		HeaderAlignment: Default.HeaderAlignment,

		RowLine:         Default.RowLine,
		ColumnSeparator: Default.ColumnSeparator,
		NewLine:         Default.NewLine,
		CenterSeparator: Default.CenterSeparator,

		DefaultAlignment: Default.DefaultAlignment,
		NumbersAlignment: Default.NumbersAlignment,

		RowLengthTitle: Default.RowLengthTitle,
	}
}

// Print calls and returns the result of the `Default.Print`,
// take a look at the `Printer#Print` function for details.
func Print(w io.Writer, v interface{}, filters ...interface{}) int {
	return New(w).Print(v, filters...)
}

// RE_TODO:

// type cursor struct {
// 	headerIndex int
// }

var emptyStruct = struct{}{}

func collect(v reflect.Value, filters []interface{}) (headers []string, rows [][]string, numbersColsPosition []int) {
	v = indirectValue(v)
	kind := v.Kind()

	if kind == reflect.String {
		// no headers, but rows.
		rows = append(rows)
		return
	}

	if kind == reflect.Slice {
		var tmp = make(map[reflect.Type]struct{})

		for i, n := 0, v.Len(); i < n; i++ {
			item := indirectValue(v.Index(i))

			f := MakeFilters(item, filters...)
			if !CanAcceptRow(item, f) {
				continue
			}

			if item.Kind() != reflect.Struct {
				// if not struct, don't search its fields, just put a row as it's.
				c, r := extractCells(i, emptyHeader, indirectValue(item))
				rows = append(rows, r)
				numbersColsPosition = append(numbersColsPosition, c...)
				continue
			}

			c, r := GetRow(item)
			numbersColsPosition = append(numbersColsPosition, c...)

			itemTyp := item.Type()
			if _, ok := tmp[itemTyp]; !ok {
				// make headers once per type.
				tmp[itemTyp] = emptyStruct
				hs := extractHeaders(itemTyp)
				if len(hs) == 0 {
					continue
				}
				for _, h := range hs {
					headers = append(headers, h.Name)
				}
			}

			rows = append(rows, r)
		}

		return
	} else if kind == reflect.Struct {
		hs := extractHeaders(v.Type())
		if len(hs) == 0 {
			return
		}

		for _, h := range hs {
			headers = append(headers, h.Name)
		}

		f := MakeFilters(v, filters...)
		if !CanAcceptRow(v, f) {
			return
		}

		c, r := GetRow(v)
		rows = append(rows, r)
		numbersColsPosition = c
	}

	return
}

var emptyHeaders []string

func (p *Printer) acquireTable() *tablewriter.Table {
	if v := p.pool.Get(); v != nil {
		return v.(*tablewriter.Table)
	}

	table := tablewriter.NewWriter(p.out)

	table.SetAlignment(int(p.DefaultAlignment))
	table.SetAutoFormatHeaders(p.AutoFormatHeaders)
	table.SetAutoWrapText(p.AutoWrapText)
	table.SetBorders(tablewriter.Border{Top: p.BorderTop, Left: p.BorderLeft, Right: p.BorderRight, Bottom: p.BorderBottom})
	table.SetHeaderLine(p.HeaderLine)
	table.SetHeaderAlignment(int(p.HeaderAlignment))
	table.SetRowLine(p.RowLine)
	table.SetColumnSeparator(p.ColumnSeparator)
	table.SetNewLine(p.NewLine)
	table.SetCenterSeparator(p.CenterSeparator)

	return table
}

func (p *Printer) releaseTable(table *tablewriter.Table) {
	table.ClearRows()
	table.SetHeader(emptyHeaders)

	p.pool.Put(table)
}

func (p *Printer) render(headers []string, rows [][]string, numbersColsPosition []int) int {
	table := p.acquireTable()
	defer p.releaseTable(table)

	if len(headers) == 0 {
		return 0
	}

	if p.RowLengthTitle != nil && p.RowLengthTitle(len(rows)) {
		headers[0] = fmt.Sprintf("%s (%d) ", headers[0], len(rows))
	}

	table.SetHeader(headers)
	table.AppendBulk(rows)

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

// Print usage:
// Print(writer, tt, func(t MyStruct) bool { /* or any type, depends on the type(s) of the "tt" */
// 	return t.Visibility != "hidden"
// })
//
// Returns the number of rows finally printed.
func (p *Printer) Print(in interface{}, filters ...interface{}) int {
	h, r, c := collect(reflect.ValueOf(in), filters)
	return p.render(h, r, c)
}

func PrintHeadList(w io.Writer, list interface{}, header string, filters ...interface{}) int {
	return New(w).PrintHeadList(list, header, filters...)
}

var emptyHeader Header

func (p *Printer) PrintHeadList(list interface{}, header string, filters ...interface{}) int {
	items := indirectValue(reflect.ValueOf(list))
	if items.Kind() != reflect.Slice {
		return 0
	}

	var (
		rows                [][]string
		numbersColsPosition []int
	)

	for i, n := 0, items.Len(); i < n; i++ {
		item := items.Index(i)
		c, r := extractCells(i, emptyHeader, indirectValue(item))
		rows = append(rows, r)
		numbersColsPosition = append(numbersColsPosition, c...)
	}

	headers := []string{header}
	return p.render(headers, rows, numbersColsPosition)
}

func PrintMap(w io.Writer, m interface{}, filters ...interface{}) int {
	return New(w).PrintMap(m, filters...)
}

// DEPRECATED by `map.go`, save the current work and go back to papers.
func (p *Printer) PrintMap(m interface{}, filters ...interface{}) int {
	v := indirectValue(reflect.ValueOf(m))
	if v.Kind() != reflect.Map {
		return 0
	}

	keys := v.MapKeys()

	if len(keys) == 0 {
		return 0
	}

	if keys[0].Kind() != reflect.String {
		return -1 // all keys should be as string, they are the header(s).
	}

	var (
		headers             []string
		rows                [][]string
		numbersColsPosition []int
	)

	for i, key := range keys {
		header := key.String()
		hasAlready := false
		for _, h := range headers {
			if h == header {
				hasAlready = true
			}
		}
		if !hasAlready {
			headers = append(headers, header)
		}

		item := v.MapIndex(key)

		if item.Kind() == reflect.Slice {
			c, r := extractCells(i, emptyHeader, item)
			rows = append(rows, r)
			numbersColsPosition = append(numbersColsPosition, c...)
		}
	}

	return p.render(headers, rows, numbersColsPosition)
}
