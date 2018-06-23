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

type structParser struct {
	TagsOnly bool
}

func (p *structParser) Parse(v reflect.Value, filters []RowFilter) ([]string, [][]string, []int) {
	if !CanAcceptRow(v, filters) {
		return nil, nil, nil
	}

	row, nums := p.ParseRow(v)

	return p.ParseHeaders(v), [][]string{row}, nums
}

func (p *structParser) ParseHeaders(v reflect.Value) []string {
	hs := extractHeadersFromStruct(v.Type(), true)
	if len(hs) == 0 {
		return nil
	}

	headers := make([]string, len(hs))
	for idx := range hs {
		headers[idx] = hs[idx].Name
	}

	return headers
}

func (p *structParser) ParseRow(v reflect.Value) ([]string, []int) {
	return getRowFromStruct(v, p.TagsOnly)
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

func extractHeaderFromStructField(f reflect.StructField, pos int, tagsOnly bool) (header StructHeader, ok bool) {
	headerTag := f.Tag.Get(HeaderTag)
	if headerTag == "" && tagsOnly {
		return emptyHeader, false
	}

	// embedded structs are acting like headers appended to the existing(s).
	if f.Type.Kind() == reflect.Struct && headerTag == InlineHeaderTag {
		return extractHeaderFromStructField(f, pos, tagsOnly)
	} else if headerTag != "" {
		if header, ok := extractHeaderFromTag(headerTag); ok {
			header.Position = pos
			return header, true
		}

	} else if !tagsOnly {
		return StructHeader{
			Position: pos,
			Name:     f.Name,
		}, true
	}

	return emptyHeader, false
}

func extractHeadersFromStruct(typ reflect.Type, tagsOnly bool) (headers []StructHeader) {
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
		header, ok := extractHeaderFromStructField(f, i, tagsOnly)
		if ok {
			headers = append(headers, header)
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
func getRowFromStruct(v reflect.Value, tagsOnly bool) (cells []string, rightCells []int) {
	typ := v.Type()
	j := 0
	for i, n := 0, typ.NumField(); i < n; i++ {
		header, ok := extractHeaderFromStructField(typ.Field(i), i, tagsOnly)
		if !ok {
			continue
		}

		fieldValue := indirectValue(v.Field(i))
		c, r := extractCells(j, header, fieldValue, tagsOnly)
		rightCells = append(rightCells, c...)
		cells = append(cells, r...)
		j++
	}

	return
}
