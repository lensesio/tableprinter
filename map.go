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
	Parse(reflect.Value) (headers []string, rows [][]string, numbers []int)
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

func (r *MapParser) Parse(v reflect.Value) ([]string, [][]string, []int) {
	keys := r.Keys(v)
	headers := r.ParseHeaders(v, keys)
	rows, numbers := r.ParseRows(v, keys)

	return headers, rows, numbers
}

func (r *MapParser) Keys(v reflect.Value) []reflect.Value {
	return v.MapKeys()
}

func (r *MapParser) ParseRows(v reflect.Value, keys []reflect.Value) ([][]string, []int) {
	// cursors := make(map[int]int) // key = map's key index(although maps don't keep order), value = current index of elements inside the map.

	maxLength := maxMapElemLength(v, keys)
	rows := make([][]string, maxLength, maxLength)
	// depends on the header size, this is for the entire col aligment but
	// we can't do that on `GetHeaders` because its values depends on the rows[index] value's type to the table.
	numbers := make([]int, 0)

	// c := 0
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
			logger.Debugf("%s:", tryString(key))
		}

		elem := v.MapIndex(key)
		if elem.Kind() != reflect.Slice {
			panic("TODO")
		}

		for i, n := 0, elem.Len(); i < n; i++ {
			// cursors[c] = i
			item := elem.Index(i)
			a, row := extractCells(i, emptyHeader, item)

			if r.Debug {
				logger.Debugf("%s%s", strings.Repeat(" ", len(tryString(key))), tryString(item))
			}

			rows[i] = append(rows[i], row...)
			numbers = append(numbers, a...)
		}

		//	c++
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

		if header := tryString(key); header != "" {
			headers = append(headers, header)
		}
	}

	return
}

func maxMapElemLength(v reflect.Value, keys []reflect.Value) (max int) {
	for _, key := range keys {
		elem := v.MapIndex(key)

		if current := elem.Len(); current > max {
			max = current
		}
	}

	return
}

func tryString(key reflect.Value) string {
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
