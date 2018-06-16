package tableprinter

import (
	"reflect"
	"strings"
)

// Cell contains the necessary information about a cell, useful for its presentation
// such as alignment, alternative value if main is empty, if the row should print the number of elements inside a list or if the column should be formated as number.
type Cell struct {
	ValueAsNumber    bool
	ValueAsCountable bool
	AlternativeValue string
	Header           string
}

// GetCell accepts a struct's field and returns its nessecary information for its presentation.
// It returns false if the particular struct's field does not contain the `HeaderTag`.
func GetCell(f reflect.StructField) (Cell, bool) {
	var c Cell

	headerTagValue := f.Tag.Get(HeaderTag)
	if headerTagValue == "" {
		return c, false
	}

	headerValues := strings.Split(headerTagValue, ",")
	switch len(headerValues) {
	case 0, 1:
		c.Header = headerTagValue
		break
	default:
		c.Header = headerValues[0]
		headerValues = headerValues[1:] /* except the first which should be the header value */
		for _, hv := range headerValues {
			switch hv {
			case NumberHeaderTag:
				c.ValueAsNumber = true
				break
			case CountHeaderTag:
				c.ValueAsCountable = true
				break
			default:
				c.AlternativeValue = hv
			}
		}
	}

	return c, true
}
