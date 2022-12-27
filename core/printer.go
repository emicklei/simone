package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/emicklei/simone/api"
)

var printer = &Printer{
	registry: map[reflect.Type]api.PrintFunc{},
}

func RegisterPrinter(v any, p api.PrintFunc) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	printer.registry[rt] = p
}

type Printer struct {
	registry map[reflect.Type]api.PrintFunc
}

func Print(v any) string {
	b := new(strings.Builder)
	printOn(v, b)
	return b.String()
}

func printOn(v any, b *strings.Builder) {
	if v == nil {
		b.WriteString("null")
		return
	}
	// check for struct
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	if rt.Kind() == reflect.Struct {
		if pf, ok := printer.registry[rt]; ok {
			pf(v, b)
			return
		}
		printStruct(b, rt, rv)
		return
	}
	// fallback to standard JSON encoder
	data, _ := json.Marshal(v)
	b.WriteString(string(data))
}

func printStruct(b *strings.Builder, rt reflect.Type, rv reflect.Value) {
	if !rv.IsValid() {
		return
	}
	b.WriteRune('{')
	comma := false
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.IsExported() {
			if comma {
				b.WriteString(", ")
			}
			b.WriteRune('\'')
			fmt.Fprintf(b, f.Name)
			b.WriteRune('\'')
			fv := rv.Field(i)
			b.WriteRune(':')
			if fv.CanInterface() {
				printOn(fv.Interface(), b)
			}
			comma = true
		}
	}
	b.WriteRune('}')
}
