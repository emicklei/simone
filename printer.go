package simone

import (
	"fmt"
	"reflect"
	"strconv"
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
	if v == nil {
		return "null"
	}
	b := new(strings.Builder)
	printOn(v, b)
	return b.String()
}

func printOn(v any, b *strings.Builder) {
	switch val := v.(type) {
	case string:
		b.WriteRune('\'')
		b.WriteString(val)
		b.WriteRune('\'')
	case int:
		b.WriteString(strconv.Itoa(val))
	case bool:
		if val {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case *bool:
		printOn(*val, b)
	default:
		rt := reflect.TypeOf(v)
		rv := reflect.ValueOf(v)
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			rv = rv.Elem()
		}
		if rt.Kind() == reflect.Struct {
			if pf, ok := printer.registry[rt]; ok {
				pf(v, b)
			} else {
				printStruct(b, rt, rv)
			}
		} else {
			fmt.Fprintf(b, "%v", v)
		}
	}
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
