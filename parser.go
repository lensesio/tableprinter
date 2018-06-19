package tableprinter

import (
	"reflect"
	"strings"
	"sync"
)

// Headers are being cached from the root-level structure to print out.
// They can be customized for custom head titles.
//
// Header can also contain the necessary information about its values, useful for its presentation
// such as alignment, alternative value if main is empty, if the row should print the number of elements inside a list or if the column should be formated as number.
var (
	Headers map[reflect.Type][]Header // type is the root struct.
	mu      sync.RWMutex
)

// Header contains the name of the header extracted from the struct's `HeaderTag` field tag.
type Header struct {
	Name string
	// Position is the horizontal position (start from zero) of the header.
	Position int

	ValueAsNumber    bool
	ValueAsCountable bool
	AlternativeValue string
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

func extractHeaders(typ reflect.Type) (headers []Header) {
	typ = indirectType(typ)
	if typ.Kind() != reflect.Struct {
		return
	}

	// search cache.
	mu.RLock()
	if cached, has := Headers[typ]; has {
		return cached
	}
	mu.RUnlock()

	for i, n := 0, typ.NumField(); i < n; i++ {
		f := typ.Field(i)

		headerTag := f.Tag.Get(HeaderTag)
		// embedded structs are acting like headers appended to the existing(s).
		if f.Type.Kind() == reflect.Struct && headerTag == InlineHeaderTag {
			headers = append(headers, extractHeaders(f.Type)...)
		} else if headerTag != "" {
			if header, ok := extractHeader(headerTag); ok {
				header.Position = i
				headers = append(headers, header)
			}
		}
	}

	if len(headers) > 0 {
		// insert to cache if it's valid table.
		mu.Lock()
		Headers[typ] = headers
		mu.Unlock()
	}

	return headers
}

func extractHeader(headerTag string) (header Header, ok bool) {
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
