package tableprinter

import (
	"reflect"
	"testing"
	"time"
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

	if expected, got := reflect.ValueOf(changedUnquote).FieldByName("HeaderFieldSet").Interface().(string), sample.HeaderFieldSet; expected != got {
		t.Fatalf("[changedUnquote] expected the field value of 'HeaderFieldSet' to be: '%s' but got: '%s'", expected, got)
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

func TestExtractTimestampHeaderTag(t *testing.T) {
	tests := []struct {
		tag      string
		is       bool
		expected TimestampHeaderTagValue
	}{
		{"timestm", false, TimestampHeaderTagValue{}},
		{"timestamp(ms", true, TimestampHeaderTagValue{}}, // it can't parse the args but it's timestamp.
		{"timestamp(ms|utc)", true, TimestampHeaderTagValue{FromMilliseconds: true, UTC: true, Format: TimestampFormatRFC822ZHeaderTag /*default*/, Local: false, Human: false}},
		{"timestamp(ms|utc|RFC822Z)", true, TimestampHeaderTagValue{FromMilliseconds: true, UTC: true, Format: time.RFC822Z, Local: false, Human: false}},
		{"timestamp(ms|local|UnixDate)", true, TimestampHeaderTagValue{FromMilliseconds: true, Local: true, Format: time.UnixDate, UTC: false, Human: false}},
		{"timestamp(RubyDate)", true, TimestampHeaderTagValue{Format: time.RubyDate, FromMilliseconds: false, Local: false, UTC: false, Human: false}},
		{"timestamp(RubyDate|utc|ms)", true, TimestampHeaderTagValue{FromMilliseconds: true, UTC: true, Format: time.RubyDate, Local: false, Human: false}},
		{"timestamp(RubyDate|local|ms)", true, TimestampHeaderTagValue{FromMilliseconds: true, Local: true, Format: time.RubyDate, UTC: false, Human: false}},
		// custom format and test if the last argument overrides the prev:
		{"timestamp(RubyDate|local|ms|02 Jan 06 15:04)", true, TimestampHeaderTagValue{FromMilliseconds: true, Local: true, Format: "02 Jan 06 15:04", UTC: false, Human: false}},
	}

	for i, tt := range tests {
		v, ok := extractTimestampHeader(tt.tag)
		if tt.is && !ok {
			t.Fatalf("[%d: '%s'] expected to be a valid timestamp header tag", i, tt.tag)
		} else if !tt.is && ok {
			t.Fatalf("[%d: '%s'] expected to be an invalid timestamp header tag but extracted as valid one", i, tt.tag)
		}

		if !reflect.DeepEqual(tt.expected, v) {
			t.Fatalf("[%d: '%s'] expected the header tag value to be: %#+v but got: %#+v", i, tt.tag, tt.expected, v)
		}
	}
}
