package tableprinter

import (
	"reflect"
	"testing"
)

func TestJSONWithInvalidNoBytesOrStringValueAndEmpties(t *testing.T) {
	type sample struct {
		HeaderField   string `header:"headervalue1"`
		MultiTagField string `json:"jsonvalue1" header:"headervalue2" xml:"xmlvalue1"`
	}

	var tt []sample

	v := indirectValue(reflect.ValueOf(tt))

	// invalid type, the JSONParse accepts []string or string only.
	h, r, nums := JSONParser.Parse(v, nil)

	if expected, got := 0, len(h); expected != got {
		t.Fatalf("the length of the headers were expected %d headers but got %d", expected, got)
	}

	if expected, got := 0, nums; expected != 0 {
		t.Fatalf("expected %d number-type fields but got %d", expected, got)
	}

	if expected, got := 0, len(r); expected != got {
		t.Fatalf("expected %d rows but got %d", expected, got)
	}

	// these should not panic at least.
	var sample2 = []byte{0x00, 0x00}
	_, _, _ = JSONParser.Parse(indirectValue(reflect.ValueOf(sample2)), nil)
	var sample3 []byte
	_, _, _ = JSONParser.Parse(indirectValue(reflect.ValueOf(sample3)), nil)
	_, _, _ = JSONParser.Parse(reflect.ValueOf(nil), nil)
}
