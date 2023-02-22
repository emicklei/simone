package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"

	"github.com/emicklei/simone/api"
)

type PlainText string

var printer = &Printer{
	registry: map[reflect.Type]api.PrintFunc{},
}

func RegisterPrinter(v any, p api.PrintFunc) {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	log.Println("registered custom printer for", rt.PkgPath()+"."+rt.Name())
	printer.registry[rt] = p
}

type Printer struct {
	registry map[reflect.Type]api.PrintFunc
}

func Print(v any) string {
	if v == NoOutputValue || v == NoOutputValueString {
		return NoOutputValueString
	}
	b := new(strings.Builder)
	printOn(v, b)
	return b.String()
}

// assume plugin is pointer type
var pluginType = reflect.TypeOf((*api.Plugin)(nil)).Elem()

func printSliceOn(v any, rt reflect.Type, b *strings.Builder) {
	rv := reflect.ValueOf(v)
	b.WriteRune('[')
	for i := 0; i < rv.Len(); i++ {
		if i > 0 {
			b.WriteString(", ")
		}
		each := rv.Index(i)
		printOn(each.Interface(), b)
	}
	b.WriteRune(']')
}

func printDefaultOn(v any, b *strings.Builder) {
	// fallback to standard JSON encoder
	data, _ := json.Marshal(v)
	b.WriteString(string(data))
}

func printOn(v any, b *strings.Builder) {
	if v == nil {
		b.WriteString("null")
		return
	}
	// PlainText is to communicate non-JSON object results
	if p, ok := v.(PlainText); ok {
		io.WriteString(b, string(p))
		return
	}
	// check for struct
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Slice {
		printSliceOn(v, rt, b)
		return
	}
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
	printDefaultOn(v, b)
}

func printStruct(b *strings.Builder, rt reflect.Type, rv reflect.Value) {
	if !rv.IsValid() {
		b.WriteString("null")
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
			b.WriteRune('"')
			fmt.Fprintf(b, f.Name)
			b.WriteRune('"')
			fv := rv.Field(i)
			b.WriteRune(':')
			if fv.CanInterface() {
				if fv.Kind() == reflect.Pointer {
					fv = fv.Elem()
				}
				if fv.IsValid() {
					printOn(fv.Interface(), b)
				}
			}
			comma = true
		}
	}
	b.WriteRune('}')
}
