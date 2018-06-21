package tableprinter

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kataras/golog"
)

/* This is a design proposal, differs from the rest of the implementation, doesn't work as one yet. */

type Parser interface {
	// Why not `ParseRows` and `ParseHeaders`?
	// Because type map has not a specific order, order can change at different runtimes,
	// so we must keep record on the keys order the first time we fetche them (=> see `MapParser#ParseRows`, `MapParser#ParseHeaders`).
	Parse(reflect.Value, []RowFilter) (headers []string, rows [][]string, numbers []int)
}

var (
	logger = golog.New().SetOutput(os.Stdout).SetTimeFormat("").SetLevel("debug")

	mapParser = new(MapParser)
)

func whichParser(typ reflect.Type) Parser {
	switch typ.Kind() {
	case reflect.Map:
		return mapParser
	default:
		// TODO:...
		return nil
	}
}

// Should we have a single parser value its specific types and give input arguments to the funcs, like "keys"
// or is better to initialize a new parser on each output, so it can be used as a cache?
type MapParser struct {
	Debug bool
}

func (r *MapParser) Parse(v reflect.Value, filters []RowFilter) ([]string, [][]string, []int) {
	keys := r.Keys(v)
	headers := r.ParseHeaders(v, keys)
	rows, numbers := r.ParseRows(v, keys, filters)

	return headers, rows, numbers
}

func (r *MapParser) Keys(v reflect.Value) []reflect.Value {
	return v.MapKeys()
}

func (r *MapParser) ParseRows(v reflect.Value, keys []reflect.Value, filters []RowFilter) ([][]string, []int) {
	// cursors := make(map[int]int) // key = map's key index(although maps don't keep order), value = current index of elements inside the map.

	maxLength := maxMapElemLength(v, keys)

	rows := make([][]string, maxLength, maxLength)
	// depends on the header size, this is for the entire col aligment but
	// we can't do that on `GetHeaders` because its values depends on the rows[index] value's type to the table.
	numbers := make([]int, 0)

	for _, key := range keys {
		// Debug for output:
		/*
			[DBUG] Sellers:
			[DBUG]        Georgios Callas
			[DBUG]        Ioannis Christou
			[DBUG] Consumers:
			[DBUG]          Dimitrios Dellis
			[DBUG]          Nikolaos Doukas
		*/
		if r.Debug {
			logger.Debugf("%s:", stringValue(key))
		}

		elem := v.MapIndex(key)
		if elem.Kind() != reflect.Slice {
			if !CanAcceptRow(elem, filters) {
				continue
			}

			a, row := extractCells(0, emptyHeader, elem)
			if len(row) == 0 {
				continue
			}
			if r.Debug {
				logger.Debugf("%s%s", strings.Repeat(" ", len(stringValue(key))), stringValue(elem))
			}
			if cap(rows) == 0 {
				rows = [][]string{row}
			} else {
				rows[0] = append(rows[0], row...)
			}

			numbers = append(numbers, a...)
			continue
		}

		for i, n := 0, elem.Len(); i < n; i++ {
			// cursors[c] = i
			item := elem.Index(i)
			if !CanAcceptRow(item, filters) {
				continue
			}

			a, row := extractCells(i, emptyHeader, item)

			if len(row) == 0 {
				continue
			}

			if r.Debug {
				logger.Debugf("%s%s", strings.Repeat(" ", len(stringValue(key))), stringValue(item))
			}

			// note that we must not check when iterating, because it may be extended before or after,
			// we must somehow collect these points rx:cx and do that on the final state...
			// if shouldEmpty {
			// 	if i == n-1 && i < maxLength-1 {
			// 		rows[i] = []string{" "}
			// 	}
			// }

			rows[i] = append(rows[i], row...)
			// if i == n-1 && i < maxLength-1 {
			// 	println("shouldEmpty set to 'true'")
			// 	shouldEmpty = true
			// 	// rows[i+1] = []string{" "}
			// }

			numbers = append(numbers, a...)
		}
	}

	return rows, numbers
}

func (r *MapParser) ParseHeaders(v reflect.Value, keys []reflect.Value) (headers []string) {
	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		// support any type, even if it's declared as "interface{}" or pointer to something, we care about this "something"'s value.
		key = indirectValue(key)
		if !key.CanInterface() {
			continue
		}

		if header := stringValue(key); header != "" {
			headers = append(headers, header)
		}
	}

	return
}

func maxMapElemLength(v reflect.Value, keys []reflect.Value) (max int) {
	for _, key := range keys {
		elem := v.MapIndex(key)
		if elem.Kind() != reflect.Slice {
			continue
		}
		if current := elem.Len(); current > max {
			max = current
		}
	}

	return
}

func stringValue(key reflect.Value) string {
	if !key.CanInterface() {
		return ""
	}

	switch keyV := key.Interface().(type) {
	case string:
		return keyV
	case fmt.Stringer:
		return keyV.String()
	default:
		return ""
	}
}
