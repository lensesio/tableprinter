package tableprinter

import (
	"reflect"
	"strings"
	"sync"
)

// StructHeaders are being cached from the root-level structure to print out.
// They can be customized for custom head titles.
//
// Header can also contain the necessary information about its values, useful for its presentation
// such as alignment, alternative value if main is empty, if the row should print the number of elements inside a list or if the column should be formated as number.
var (
	StructHeaders    = make(map[reflect.Type][]StructHeader) // type is the root struct.
	structsHeadersMu sync.RWMutex
)

type structParser struct{}

func (p *structParser) Parse(v reflect.Value, filters []RowFilter) ([]string, [][]string, []int) {
	hs := extractHeadersFromStruct(v.Type())
	if len(hs) == 0 {
		return nil, nil, nil
	}

	headers := make([]string, len(hs))
	for idx := range hs {
		headers[idx] = hs[idx].Name
	}

	if !CanAcceptRow(v, filters) {
		return nil, nil, nil
	}

	nums, row := getRowFromStruct(v)

	return headers, [][]string{row}, nums
}

// StructHeader contains the name of the header extracted from the struct's `HeaderTag` field tag.
type StructHeader struct {
	Name string
	// Position is the horizontal position (start from zero) of the header.
	Position int

	ValueAsNumber    bool
	ValueAsCountable bool
	AlternativeValue string
}

func extractHeadersFromStruct(typ reflect.Type) (headers []StructHeader) {
	typ = indirectType(typ)
	if typ.Kind() != reflect.Struct {
		return
	}

	// search cache.
	structsHeadersMu.RLock()
	if cached, has := StructHeaders[typ]; has {
		structsHeadersMu.RUnlock()
		return cached
	}
	structsHeadersMu.RUnlock()

	for i, n := 0, typ.NumField(); i < n; i++ {
		f := typ.Field(i)

		headerTag := f.Tag.Get(HeaderTag)
		// embedded structs are acting like headers appended to the existing(s).
		if f.Type.Kind() == reflect.Struct && headerTag == InlineHeaderTag {
			headers = append(headers, extractHeadersFromStruct(f.Type)...)
		} else if headerTag != "" {
			if header, ok := extractHeaderFromTag(headerTag); ok {
				header.Position = i
				headers = append(headers, header)
			}
		}
	}

	if len(headers) > 0 {
		// insert to cache if it's valid table.
		structsHeadersMu.Lock()
		StructHeaders[typ] = headers
		structsHeadersMu.Unlock()
	}

	return headers
}

func extractHeaderFromTag(headerTag string) (header StructHeader, ok bool) {
	if headerTag == "" {
		return
	}
	ok = true

	parts := strings.Split(headerTag, ",")

	// header name is the first part.
	header.Name = parts[0]

	if len(parts) > 1 {
		for _, hv := range parts[1:] /* except the first part ofc which should be the header value */ {
			switch hv {
			case NumberHeaderTag:
				header.ValueAsNumber = true
				break
			case CountHeaderTag:
				header.ValueAsCountable = true
				break
			default:
				header.AlternativeValue = hv
			}
		}
	}

	return
}

// getRowFromStruct returns the positions of the cells that should be aligned to the right
// and the list of cells(= the values based on the cell's description) based on the "in" value.
func getRowFromStruct(v reflect.Value) (rightCells []int, cells []string) {
	typ := v.Type()
	j := 0
	for i, n := 0, typ.NumField(); i < n; i++ {
		header, ok := extractHeaderFromTag(typ.Field(i).Tag.Get(HeaderTag))
		if !ok {
			continue
		}

		fieldValue := indirectValue(v.Field(i))
		c, r := extractCells(j, header, fieldValue)
		rightCells = append(rightCells, c...)
		cells = append(cells, r...)
		j++
	}

	return
}
