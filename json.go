package tableprinter

import (
	"encoding/json"
	"reflect"
)

type jsonParser struct{}

func (p *jsonParser) Parse(v reflect.Value, filters []RowFilter) (headers []string, rows [][]string, nums []int) {
	var b []byte

	if kind := v.Kind(); kind == reflect.Slice {
		b = v.Bytes()
	} else if kind == reflect.String {
		b = []byte(v.String())
	} else {
		return
	}

	var in interface{} // or map[string]interface{}
	if err := json.Unmarshal(b, &in); err != nil {
		return
	}

	inValue := indirectValue(reflect.ValueOf(in))
	return WhichParser(inValue.Type()).Parse(inValue, filters)
}
