package simone

import (
	"encoding/json"
	"reflect"
	"strings"
)

var printer = &Printer{
	registry: map[reflect.Type]PrintFunc{},
}

func RegisterPrinter(v any, p PrintFunc) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	printer.registry[rt] = p
}

type Printer struct {
	registry map[reflect.Type]PrintFunc
}

func Print(v any) string {
	if v != nil {
		// first check for custom printer
		rt := reflect.TypeOf(v)
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
		}
		if rt.Kind() == reflect.Struct {
			if pf, ok := printer.registry[rt]; ok {
				b := new(strings.Builder)
				pf(v, b)
				return b.String()
			}
		}
	}
	// fallback to standard JSON encoder
	data, _ := json.Marshal(v)
	return string(data)
}
