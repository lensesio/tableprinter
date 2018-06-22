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

func TestMapEmpties(t *testing.T) {
	persons := map[string][]person{
		"Access 1": []person{{"Georgios", "Callas"},
			{"Ioannis", "Christou"}},
		"Access 2": []person{
			{"Dimitrios", "Dellis"}},
		"Access 3": []person{{"Giannhs", "Papadopoulos"},
			{"Giwrgos", "Papadopoulos"},
			{"Oresths", "Papadopoulos"}},
	}

	_, rows, _ := mapParser.Parse(reflect.ValueOf(persons), nil)
	if expected, got := len(persons), len(rows); expected != got {
		t.Fatalf("expected %d rows but got %d", expected, got)
	}

	/* Remember: This can be different from runtime to runtime, maps are not always have the same key order(= our headers),
	 so this test can fail sometimes because it checks the exact pos of those empties (prev tests are written to adjust that behavior, they should always SUCC)

		  ACCESS 1              ACCESS 2                  ACCESS 3
		 --------------------  ------------------        ----------------------
	[0]   Georgios Callas[0]    Dimitrios Dellis[0:1]   Giannhs Papadopoulos[0:2]
	[1]   Ioannis Christou[0:1] EMPTY [1:1]             Giwrgos Papadopoulos[1:2]
	[2]	  EMPTY [2:0]           EMPTY [2:1]             Oresths Papadopoulos[2]
	*/

	var (
		space = " "

		empties = map[int][]int{
			2: []int{0, 1},
			1: []int{1},
		}

		someNotEmpties = map[int][]int{
			2: []int{2},
		}
	)

	for idx, list := range empties {
		for _, e := range list {
			if got := rows[idx][e]; got != space {
				t.Fatalf("expected %d:%d to have space but got: %s", idx, e, got)
			}
		}
	}

	for idx, list := range someNotEmpties {
		for _, e := range list {
			if got := rows[idx][e]; got == space {
				t.Fatalf("expected %d:%d to have filled value but got space", idx, e)
			}
		}
	}
}
