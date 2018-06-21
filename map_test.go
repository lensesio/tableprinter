package tableprinter

import (
	"reflect"
	"testing"
)

type person struct {
	FirstName string
	LastName  string
}

func (p person) String() string {
	return p.FirstName + " " + p.LastName
}

func TestMapParse(t *testing.T) {
	t.Parallel()

	tt := map[string][]person{
		"Sellers": []person{{"Georgios", "Callas"},
			{"Ioannis", "Christou"}},
		"Consumers": []person{
			{"Dimitrios", "Dellis"},
			{"Nikolaos", "Doukas"}},
	}

	var (
		expectedHeaders = []string{"Sellers", "Consumers"}
		// the order may differs, remember map doesn't keep its order, so make a check if row contains these values.
		expectedRowPart1 = []string{"Georgios Callas", "Dimitrios Dellis"}
		expectedRowPart2 = []string{"Ioannis Christou", "Nikolaos Doukas"}
		expectedRows     = [][]string{expectedRowPart1, expectedRowPart2}
	)

	v := reflect.ValueOf(tt)
	headers, rows, _ := mapParser.Parse(v, nil)

	// check the length.
	if expected, got := len(expectedHeaders), len(headers); expected != got {
		t.Fatalf("expected length of headers: %d but got: %d", expected, got)
	}
	// we should not care about the order of headers on maps, we just check if all expected headers are there.
	if !((headers[0] == expectedHeaders[0] || headers[0] == expectedHeaders[1]) && (headers[1] == expectedHeaders[0] || headers[1] == expectedHeaders[1])) {
		t.Fatalf("expected headers: %v but got: %v", expectedHeaders, headers)
	}

	// check the length.
	if expected, got := len(expectedRows), len(rows); expected != got {
		t.Fatalf("expected length of rows: %d but got: %d", expected, got)
	}

	// we care if the correct header contains the correct columns.
	if !((rows[0][0] == expectedRows[0][0] || rows[0][1] == expectedRows[0][0]) && (rows[1][0] == expectedRows[1][0] || rows[1][1] == expectedRows[1][0])) {
		t.Fatalf("expected rows: %v but got: %v", expectedRows, rows)
	}
}

func TestMapParseSingle(t *testing.T) {
	t.Parallel()

	tt := map[string]person{
		"Seller":   person{"Georgios", "Callas"},
		"Consumer": person{"Dimitrios", "Dellis"},
	}

	var (
		expectedHeaders = []string{"Seller", "Consumer"}
		expectedRows    = [][]string{[]string{"Georgios Callas", "Dimitrios Dellis"}}
	)

	v := reflect.ValueOf(tt)
	headers, rows, _ := mapParser.Parse(v, nil)

	// check the length.
	if expected, got := len(expectedHeaders), len(headers); expected != got {
		t.Fatalf("expected length of headers: %d but got: %d", expected, got)
	}
	// we should not care about the order of headers on maps, we just check if all expected headers are there.
	if !((headers[0] == expectedHeaders[0] || headers[0] == expectedHeaders[1]) && (headers[1] == expectedHeaders[0] || headers[1] == expectedHeaders[1])) {
		t.Fatalf("expected headers: %v but got: %v", expectedHeaders, headers)
	}

	// check the length.
	if expected, got := len(expectedRows), len(rows); expected != got {
		t.Fatalf("expected length of rows: %d but got: %d", expected, got)
	}

	// we care if the correct header contains the correct columns.
	if !(rows[0][0] == expectedRows[0][0] || rows[0][1] == expectedRows[0][0]) {
		t.Fatalf("expected rows: %v but got: %v", expectedRows, rows)
	}
}

func TestMapEmptiesTODO(t *testing.T) {
	tt := map[string][]person{
		"Sellers": []person{
			{"Georgios", "Callas"},
			{"Ioannis", "Christou"}},
		"Consumers": []person{
			{"Dimitrios", "Dellis"},
			{"Nikolaos", "Doukas"},
			{"Ext FName", "Ext LName"}},
	}

	v := reflect.ValueOf(tt)
	_, rows, _ := mapParser.Parse(v, nil)

	if len(rows) != 3 {
		t.Fatalf("all three rows should be printed")
	}

	if rows[2][0] != " " {
		t.Fatalf("expected 2:0 to have space, this can be fixed easly but i've not found the complete solution yet because we may have more than two and many empties on different positions, left*c or right*c")
	}
}
