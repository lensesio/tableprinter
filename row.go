package tableprinter

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
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
		if filterTyp.Kind() != reflect.Func || filterTyp.NumIn() != 1 /* not receiver */ || filterTyp.In(0) != in.Type() {
			continue
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
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
			header.ValueAsNumber = true
			s = fmt.Sprintf("%d", vi)
			break
		case reflect.Float32, reflect.Float64:
			s = fmt.Sprintf("%.2f", vi)
			rightCells = append(rightCells, pos)
			break
		case reflect.Bool:
			if vi.(bool) {
				s = "Yes"
			} else {
				s = "No"
			}
			break
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
						if hasMore := fieldSliceIdx+1 > fieldSliceLen; hasMore {
							s += ", "
						}
					}
				}
			}
			break
		default:
			if viTyp := reflect.TypeOf(vi); viTyp.Kind() == reflect.Struct {
				rr, rightEmbeddedSlices := getRowFromStruct(reflect.ValueOf(vi), whenStructTagsOnly)
				if len(rr) > 0 {
					cells = append(cells, rr...)
					for range rightEmbeddedSlices {
						rightCells = append(rightCells, pos)
						pos++
					}

					return
				}
			}

			s = fmt.Sprintf("%v", vi)
		}

		if header.ValueAsNumber {
			sInt64, err := strconv.ParseInt(fmt.Sprintf("%v", s), 10, 64)
			if err != nil || sInt64 == 0 {
				s = header.AlternativeValue
				if s == "" {
					s = "0"
				}
			} else {
				s = nearestThousandFormat(float64(sInt64))
			}

			rightCells = append(rightCells, pos)
		}

		if s == "" {
			s = header.AlternativeValue
		}

		cells = append(cells, s)
	}

	return
}
