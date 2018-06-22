package tableprinter

import (
	"reflect"
)

type Parser interface {
	// Why not `ParseRows` and `ParseHeaders`?
	// Because type map has not a specific order, order can change at different runtimes,
	// so we must keep record on the keys order the first time we fetche them (=> see `MapParser#ParseRows`, `MapParser#ParseHeaders`).
	Parse(reflect.Value, []RowFilter) (headers []string, rows [][]string, numbers []int)
}

var (
	StructParser = new(structParser)
	SliceParser  = new(sliceParser)
	MapParser    = new(mapParser)
)

func whichParser(typ reflect.Type) Parser {
	switch typ.Kind() {
	case reflect.Struct:
		return StructParser
	case reflect.Slice:
		return SliceParser
	case reflect.Map:
		return MapParser
	default:
		// TODO:...
		return nil
	}
}

// like reflect.Indirect but for types and reflect.Interface types too.
func indirectType(typ reflect.Type) reflect.Type {
	if kind := typ.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		return typ.Elem()
	}

	return typ
}

// like reflect.Indirect but reflect.Interface values too.
func indirectValue(val reflect.Value) reflect.Value {
	if kind := val.Kind(); kind == reflect.Interface || kind == reflect.Ptr {
		return val.Elem()
	}

	return val
}
