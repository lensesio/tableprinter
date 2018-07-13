package tableprinter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

const (
	// HeaderTag usage: Field string `header:"Name"`
	HeaderTag = "header"
	// InlineHeaderTag usage: Embedded Struct `header:"inline"`
	InlineHeaderTag = "inline"
	// NumberHeaderTag usage: NumberButString string `header:"Age,number"`
	NumberHeaderTag = "number"
	// CountHeaderTag usage: List []any `header:"MyList,count"`
	CountHeaderTag = "count"
	// ForceTextHeaderTag usage: ID int `header:"ID,text"`
	ForceTextHeaderTag = "text"

	// TimestampHeaderTag usage: Timestamp int64 `json:"timestamp" yaml:"Timestamp" header:"At,timestamp(ms|utc|02 Jan 2006 15:04)"`
	TimestampHeaderTag = "timestamp"
	// TimestampFromMillisecondsHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms)"`
	TimestampFromMillisecondsHeaderTag = "ms"
	// TimestampAsUTCHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc)"`
	TimestampAsUTCHeaderTag = "utc"
	// TimestampAsLocalHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|local)"`
	TimestampAsLocalHeaderTag = "local"
	// TimestampFormatHumanHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|human)"`
	TimestampFormatHumanHeaderTag = "human"
	// TimestampFormatANSICHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|ANSIC)"`
	TimestampFormatANSICHeaderTag = "ANSIC"
	// TimestampFormatUnixDateCHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|UnixDate)"`
	TimestampFormatUnixDateCHeaderTag = "UnixDate"
	// TimestampFormatRubyDateHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RubyDate)"`
	TimestampFormatRubyDateHeaderTag = "RubyDate"
	// TimestampFormatRFC822HeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC822)"`
	TimestampFormatRFC822HeaderTag = "RFC822"
	// TimestampFormatRFC822ZHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC822Z)"`
	TimestampFormatRFC822ZHeaderTag = "RFC822Z"
	// TimestampFormatRFC850HeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC850)"`
	TimestampFormatRFC850HeaderTag = "RFC850"
	// TimestampFormatRFC1123HeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC1123)"`
	TimestampFormatRFC1123HeaderTag = "RFC1123"
	// TimestampFormatRFC1123ZHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC1123Z)"`
	TimestampFormatRFC1123ZHeaderTag = "RFC1123Z" // default one.
	// TimestampFormatRFC3339HeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC3339)"`
	TimestampFormatRFC3339HeaderTag = "RFC3339"
	// TimestampFormatARFC3339NanoHeaderTag usage: Timestamp int64 `header:"Start,timestamp(ms|utc|RFC3339Nano)"`
	TimestampFormatARFC3339NanoHeaderTag = "RFC3339Nano"

	// DurationHeaderTag usage: Uptime int64 `header:"Uptime,unixduration"`
	DurationHeaderTag = "unixduration"
	// DateHeaderTag usage: Start string `header:"Start,date"`, the field's value should be formatted as time.RFC3339
	DateHeaderTag = "date"
)

// RowFilter is the row's filter, accepts the reflect.Value of the custom type,
// and returns true if the particular row can be included in the final result.
type RowFilter func(reflect.Value) bool

// CanAcceptRow accepts a value of row and a set of filter
// and returns true if it can be printed, otherwise false.
// If no filters passed then it returns true.
func CanAcceptRow(in reflect.Value, filters []RowFilter) bool {
	acceptRow := true
	for _, filter := range filters {
		if filter == nil {
			continue
		}

		if !filter(in) {
			acceptRow = false
			break
		}
	}

	return acceptRow
}

var (
	rowFilters   = make(map[reflect.Type][]RowFilter)
	rowFiltersMu sync.RWMutex
)

// MakeFilters accept a value of row and generic filters and returns a set of typed `RowFilter`.
//
// Usage:
// in := reflect.ValueOf(myNewStructValue)
// filters := MakeFilters(in, func(v MyStruct) bool { return _custom logic here_ })
// if CanAcceptRow(in, filters) { _custom logic here_ }
func MakeFilters(in reflect.Value, genericFilters ...interface{}) (f []RowFilter) {
	typ := in.Type()

	rowFiltersMu.RLock()
	if cached, has := rowFilters[typ]; has {
		rowFiltersMu.RUnlock()
		return cached
	}
	rowFiltersMu.RUnlock()

	for _, filter := range genericFilters {
		filterTyp := reflect.TypeOf(filter)
		// must be a function that accepts one input argument which is the same of the "v".
		if filterTyp.Kind() != reflect.Func || filterTyp.NumIn() != 1 /* not receiver */ {
			continue
		}

		if filterInTyp := filterTyp.In(0); filterInTyp != in.Type() {
			goodElementType := false
			if in.Kind() == reflect.Slice {
				if in.Len() > 0 {
					if filterInTyp == in.Index(0).Type() {
						// the slice contains element that is the same as the filter's func, we must allow that for slices because slice parser executes that(correctly) per ELEMENT.
						goodElementType = true
					}
				}
			}
			if !goodElementType {
				continue
			}
		}

		// must be a function that returns a single boolean value.
		if filterTyp.NumOut() != 1 || filterTyp.Out(0).Kind() != reflect.Bool {
			continue
		}

		filterValue := reflect.ValueOf(filter)
		func(filterValue reflect.Value) {
			f = append(f, func(in reflect.Value) bool {
				out := filterValue.Call([]reflect.Value{in})
				return out[0].Interface().(bool)
			})
		}(filterValue)
	}

	// insert to cache, even if filters are empty.
	rowFiltersMu.Lock()
	rowFilters[typ] = f
	rowFiltersMu.Unlock()

	return
}

func extractCells(pos int, header StructHeader, v reflect.Value, whenStructTagsOnly bool) (rightCells []int, cells []string) {
	if v.CanInterface() {
		s := ""
		vi := v.Interface()

		switch v.Kind() {
		case reflect.Int64:
			if header.ValueAsTimestamp {
				n := vi.(int64)
				if n <= 0 {
					break
				}

				if header.TimestampValue.FromMilliseconds { // to seconds.
					n = n / 1000
				}

				t := time.Unix(n, 0)
				if t.IsZero() {
					break
				}

				if header.TimestampValue.UTC {
					t = t.UTC()
				} else if header.TimestampValue.Local {
					t = t.Local()
				}

				if header.TimestampValue.Human {
					s = humanize.Time(t)
				} else {
					s = t.Format(header.TimestampValue.Format)
				}
				break
			}

			if header.ValueAsDuration {
				dur := vi.(int64)
				if dur <= 0 {
					break
				}

				since := time.Now().Truncate( /*- -> it's raw duration */ time.Duration(dur))
				hour, min, sec := since.Clock()
				if hour+min == 0 {
					// if only seconds.
					s = fmt.Sprintf("%d seconds", sec)
				} else {
					s = fmt.Sprintf("%dh %dm %ds", hour, min, sec)
				}

				break
			}

			if !header.ValueAsText {
				header.ValueAsNumber = true
				rightCells = append(rightCells, pos)
			}

			s = fmt.Sprintf("%d", vi)
		// 	fallthrough
		case reflect.Int, reflect.Int16, reflect.Int32:
			if !header.ValueAsText {
				header.ValueAsNumber = true
				rightCells = append(rightCells, pos)
			}

			s = fmt.Sprintf("%d", vi)
		case reflect.Float32, reflect.Float64:
			s = fmt.Sprintf("%.2f", vi)
			rightCells = append(rightCells, pos)
		case reflect.Bool:
			if vi.(bool) {
				s = "Yes"
			} else {
				s = "No"
			}
		case reflect.Slice, reflect.Array:
			n := v.Len()
			if header.ValueAsCountable {
				s = strconv.Itoa(n)
				header.ValueAsNumber = true
			} else if n == 0 && header.AlternativeValue != "" {
				s = header.AlternativeValue
			} else {
				for fieldSliceIdx, fieldSliceLen := 0, v.Len(); fieldSliceIdx < fieldSliceLen; fieldSliceIdx++ {
					vf := v.Index(fieldSliceIdx)
					if vf.CanInterface() {
						s += fmt.Sprintf("%v", vf.Interface())
						if hasMore := fieldSliceIdx+1 < fieldSliceLen; hasMore {
							s += ", "
						}
					}
				}
			}
		case reflect.Map:
			if len(v.MapKeys()) == 0 {
				return
			}
			b, err := json.MarshalIndent(vi, " ", "  ")
			if err != nil {
				s = fmt.Sprintf("%v", vi)
			} else {
				b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
				b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
				b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)
				s = string(b)
			}

		default:
			switch t := vi.(type) {
			// Give priority to String() string functions inside the struct, if it's there then it's the whole cell string,
			// otherwise if it's struct it's the fields if TagsOnly == false, useful for dynamic maps.
			case fmt.Stringer:
				s = t.String()
			case struct{}:
				rr, rightEmbeddedSlices := getRowFromStruct(reflect.ValueOf(vi), whenStructTagsOnly)
				if len(rr) > 0 {
					cells = append(cells, rr...)
					for range rightEmbeddedSlices {
						rightCells = append(rightCells, pos)
						pos++
					}

					return
				}
			default:
				s = fmt.Sprintf("%v", vi)
			}
		}

		if header.ValueAsNumber {
			sInt64, err := strconv.ParseInt(s, 10, 64)
			if err != nil || sInt64 == 0 {
				s = header.AlternativeValue
				if s == "" {
					s = "0"
				}
			} else {
				s = nearestThousandFormat(float64(sInt64))
			}

			rightCells = append(rightCells, pos)
		} else if header.ValueAsDate {
			t, err := time.Parse(time.RFC3339, s)
			if err == nil {
				s = t.Format("2006-01-02 15:04:05")
			}
		}

		if s == "" {
			s = header.AlternativeValue
		}

		cells = append(cells, s)
	}

	return
}
