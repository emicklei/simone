package core

import (
	"encoding/json"
	"fmt"
	"io"
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

func printMethods(b *strings.Builder, rt reflect.Type) {
	it := rt
	if rt.Kind() == reflect.Pointer {
		it = rt.Elem()
	}
	fmt.Fprintf(b, "%s.%s [\n", it.PkgPath(), it.Name())
	for m := 0; m < rt.NumMethod(); m++ {
		met := rt.Method(m)
		if met.IsExported() {
			// part of Plugin interface
			if met.Name == "Init" || met.Name == "Namespace" {
				continue
			}
			printMethod(b, met)
			fmt.Fprintln(b)
		}
	}
	fmt.Fprintf(b, "]")
}

func printMethod(b *strings.Builder, met reflect.Method) {
	fmt.Fprintf(b, "\t%s(", met.Name)
	t := met.Func.Type()
	if t.Kind() != reflect.Func {
		fmt.Fprintf(b, "<not a function>:%s", t.Kind().String())
		return
	}
	// 0 = receiver
	for i := 1; i < t.NumIn(); i++ {
		if i > 1 {
			b.WriteString(", ")
		}
		b.WriteString(t.In(i).String())
	}
	b.WriteString(")")
	if numOut := t.NumOut(); numOut > 0 {
		if numOut > 1 {
			b.WriteString(" (")
		} else {
			b.WriteString(" ")
		}
		for i := 0; i < t.NumOut(); i++ {
			if i > 0 {
				b.WriteString(", ")
			}
			b.WriteString(t.Out(i).String())
		}
		if numOut > 1 {
			b.WriteString(")")
		}
	}
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
