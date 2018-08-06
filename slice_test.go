package tableprinter

import (
	"reflect"
	"testing"
)

func TestSliceEmpties(t *testing.T) {
	type sample struct {
		HeaderField   string `header:"headervalue1"`
		MultiTagField string `json:"jsonvalue1" header:"headervalue2" xml:"xmlvalue1"`
	}

	var tt []sample

	v := indirectValue(reflect.ValueOf(tt)) // like `Print` does, because we need the underline value, i.e it could be an interface{}
	parser := WhichParser(v.Type())

	// if no elements inside, then headers, rows and number-types positions should be zero.
	h, r, nums := parser.Parse(v, nil)

	if expected, got := 0, len(h); expected != got {
		t.Fatalf("the length of the headers were expected %d headers but got %d", expected, got)
	}

	if expected, got := 0, nums; expected != 0 {
		t.Fatalf("expected %d number-type fields but got %d", expected, got)
	}

	if expected, got := 0, len(r); expected != got {
		t.Fatalf("expected %d rows but got %d", expected, got)
	}
}
