package tableprinter

import (
	"reflect"
	"testing"
)

func TestSetStructHeader(t *testing.T) {
	sample := struct {
		HeaderField     string `header:"headervalue1"`
		unexportedField string `header:"we dont care"`
		MultiTagField   string `json:"jsonvalue1" header:"headervalue2" xml:"xmlvalue1"`
		HeaderFieldSet  string `header:"headervalue3"`
	}{"value1", "value2", "value3", "value4"} // it can be empty.

	// set, unquoted.
	expectedNewHeaderValueUnquote := "headervalue3_new_unquote"
	changedUnquote := SetStructHeader(sample, "HeaderFieldSet", expectedNewHeaderValueUnquote)
	typ := reflect.TypeOf(changedUnquote)
	f, ok := typ.FieldByName("HeaderFieldSet")
	if !ok {
		t.Fatalf("[changedUnquote] 'HeaderFieldSet' not found")
	}
	if expected, got := `header:"`+expectedNewHeaderValueUnquote+`"`, string(f.Tag); expected != got {
		t.Fatalf("[changedUnquote] expected the whole field tag of 'HeaderFieldSet' to be changed to '%s', but got: '%s'", expected, got)
	}

	// set, quoted.
	expectedNewHeaderValueQuote := `"headervalue3_new_quote"`
	changedQuote := SetStructHeader(sample, "HeaderFieldSet", expectedNewHeaderValueQuote)
	typ = reflect.TypeOf(changedQuote)
	f, ok = typ.FieldByName("HeaderFieldSet")
	if !ok {
		t.Fatalf("[changedQuote] 'HeaderFieldSet' not found")
	}
	if expected, got := "header:"+expectedNewHeaderValueQuote, string(f.Tag); expected != got {
		t.Fatalf("[changedQuote] expected the whole field tag of 'HeaderFieldSet' to be changed to '%s', but got: '%s'", expected, got)
	}

	t.Run("Remove", func(t *testing.T) {
		notChanged := RemoveStructHeader(sample, "unexportedField")
		if !reflect.DeepEqual(notChanged, sample) {
			t.Fatalf("[notChanged] expected the whole value to be exactly the same, original value(%#+v) should be returned instead of: %#+v", sample, notChanged)
		}

		removedHeaderField := RemoveStructHeader(sample, "HeaderField")
		typ := reflect.TypeOf(removedHeaderField)
		f, ok := typ.FieldByName("HeaderField")
		if !ok {
			t.Fatalf("[removedHeaderField] 'HeaderField' not found")
		}
		if got := string(f.Tag); got != "" {
			t.Fatalf("[removedHeaderField] expected the whole field tag of 'HeaderField' to be removed, but got: '%s'", got)
		}

		removedMultiTagField := RemoveStructHeader(sample, "MultiTagField")
		typ = reflect.TypeOf(removedMultiTagField)
		f, ok = typ.FieldByName("MultiTagField")
		if !ok {
			t.Fatalf("[removedMultiTagField] 'MultiTagField' not found")
		}
		if expected, got := `json:"jsonvalue1" xml:"xmlvalue1"`, string(f.Tag); expected != got {
			t.Fatalf("[removedMultiTagField] expected only the header tag of 'MultiTagField' to be removed('%s') but got: '%s'", expected, got)
		}
	})
}
