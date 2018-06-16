package tableprinter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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

// GetHeaders accept a type of row and returns the available headers, including the embedded headers (if any).
func GetHeaders(typ reflect.Type) (headers []string) {
	for i, n := 0, typ.NumField(); i < n; i++ {
		f := typ.Field(i)

		header := f.Tag.Get(HeaderTag)
		// embedded structs are acting like headers appended to the existing(s).
		if f.Type.Kind() == reflect.Struct && header == InlineHeaderTag {
			headers = append(headers, GetHeaders(f.Type)...)
		} else if header != "" {
			// header is the first part.
			headers = append(headers, strings.Split(header, ",")[0])
		}
	}

	return
}

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

// MakeFilters accept a value of row and generic filters and returns a set of typed `RowFilter`.
//
// Usage:
// in := reflect.ValueOf(myNewStructValue)
// filters := MakeFilters(in, func(v MyStruct) bool { return _custom logic here_ })
// if CanAcceptRow(in, filters) { _custom logic here_ }
func MakeFilters(in reflect.Value, genericFilters ...interface{}) (f []RowFilter) {
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

	return
}

// GetRow returns the positions of the cells that should be aligned to the right
// and the list of cells(= the values based on the cell's description) based on the "in" value.
func GetRow(in reflect.Value) (rightCells []int, cells []string) {
	v := reflect.Indirect(in)
	typ := v.Type()
	j := 0
	for i, n := 0, typ.NumField(); i < n; i++ {
		cell, ok := GetCell(typ.Field(i))
		if !ok {
			continue
		}

		fieldValue := reflect.Indirect(v.Field(i))

		if fieldValue.CanInterface() {
			s := ""
			vi := fieldValue.Interface()

			switch fieldValue.Kind() {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				cell.ValueAsNumber = true
				s = fmt.Sprintf("%d", vi)
				break
			case reflect.Float32, reflect.Float64:
				s = fmt.Sprintf("%.2f", vi)
				rightCells = append(rightCells, j)
				break
			case reflect.Bool:
				if vi.(bool) {
					s = "Yes"
				} else {
					s = "No"
				}
				break
			case reflect.Slice, reflect.Array:
				n := fieldValue.Len()
				if cell.ValueAsCountable {
					s = strconv.Itoa(n)
					cell.ValueAsNumber = true
				} else if n == 0 && cell.AlternativeValue != "" {
					s = cell.AlternativeValue
				} else {
					for fieldSliceIdx, fieldSliceLen := 0, fieldValue.Len(); fieldSliceIdx < fieldSliceLen; fieldSliceIdx++ {
						vf := fieldValue.Index(fieldSliceIdx)
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
					rightEmbeddedSlices, rr := GetRow(reflect.ValueOf(vi))
					if len(rr) > 0 {
						cells = append(cells, rr...)
						for range rightEmbeddedSlices {
							rightCells = append(rightCells, j)
							j++
						}

						continue
					}
				}

				s = fmt.Sprintf("%v", vi)
			}

			if cell.ValueAsNumber {
				// rightCells = append(rightCells, j)
				sInt64, err := strconv.ParseInt(fmt.Sprintf("%v", s), 10, 64)
				if err != nil || sInt64 == 0 {
					// println(err.Error())
					s = cell.AlternativeValue
					if s == "" {
						s = "0"
					}
				} else {
					s = nearestThousandFormat(float64(sInt64))
				}

				rightCells = append(rightCells, j)
			}

			if s == "" {
				s = cell.AlternativeValue
			}

			cells = append(cells, s)
			j++
		}
	}

	return
}
