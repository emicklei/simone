package core

import (
	"fmt"
	"reflect"
	"strings"
)

func printDocumentation(b *strings.Builder, v any) {
	// see if value has MethodSignatures
	rt := reflect.TypeOf(v)
	printMethods(b, rt)
}

func printMethods(b *strings.Builder, rt reflect.Type) {
	it := rt
	if rt.Kind() == reflect.Pointer {
		it = rt.Elem()
	}
	fmt.Fprintf(b, "%s.%s\n", it.PkgPath(), it.Name())
	for m := 0; m < rt.NumMethod(); m++ {
		if m > 0 {
			fmt.Fprintln(b)
		}
		met := rt.Method(m)
		if met.IsExported() {
			if isPluginMethod(met.Name) {
				continue
			}
			printMethod(b, met)
		}
	}
}

func isPluginMethod(name string) bool {
	// part of Plugin interface
	return strings.Contains("Init Namespace MethodSignatures", name)
}

func printMethod(b *strings.Builder, met reflect.Method) {
	fmt.Fprintf(b, "  %s(", met.Name)
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
