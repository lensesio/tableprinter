package tableprinter

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMakeFiltersSlice(t *testing.T) {
	var (
		persons = []person{
			{"Chris", "Doukas"},
			{"Georgios", "Callas"},
			{"Ioannis", "Christou"},
			{"Dimitrios", "Dellis"},
			{"Nikolaos", "Doukas"},
		}
		expectedLengthRows = 2
	)

	in := reflect.ValueOf(persons)

	onlyDoukasFilter := func(p person) bool {
		return p.LastName == "Doukas"
	}

	filters := MakeFilters(in, onlyDoukasFilter)
	headers, rows, nums := SliceParser.Parse(in, filters)

	if got := len(rows); got != expectedLengthRows {
		t.Fatalf("slice parser expected to return %d rows but got %d", expectedLengthRows, got)
	}

	buf := new(bytes.Buffer)
	if got := Render(buf, headers, rows, nums, false); got != expectedLengthRows {
		t.Fatalf("expected to render only %d elements containing with 'Doukas' lastname but got: %d", expectedLengthRows, got)
	}

	buf.Reset()
}
