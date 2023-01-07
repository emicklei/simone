package core

import (
	"fmt"
	"reflect"
	"strings"
)

func PrintMarkdown(v any) string {
	if v == nil {
		return "null"
	}
	b := new(strings.Builder)
	printMarkdownOn(v, b)
	return b.String()
}

func printMarkdownOn(v any, b *strings.Builder) {
	// check for struct
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	if rt.Kind() == reflect.Slice {
		printSliceMarkdownOn(rv, rt, b)
		return
	}
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	if rt.Kind() == reflect.Struct {
		printStructMarkdown(b, rt, rv)
		return
	}
	// fallback to standard JSON encoder
	printDefaultOn(v, b)
}

var maxSize = 20

func printSliceMarkdownOn(anyValue reflect.Value, sliceType reflect.Type, b *strings.Builder) {
	l := anyValue.Len()
	if l > maxSize {
		l = maxSize
	}
	for f := 0; f < l; f++ {
		each := anyValue.Index(f)
		if each.CanInterface() {
			fmt.Fprintf(b, "%d. ", f)
			printMarkdownOn(each.Interface(), b)
			b.WriteRune('\n')
		}
	}
	if l < anyValue.Len() {
		fmt.Fprintf(b, "%d. (%d more)\n", l+1, anyValue.Len()-l)
	}
}

func printStructMarkdown(b *strings.Builder, rt reflect.Type, rv reflect.Value) {
	if !rv.IsValid() {
		b.WriteString("null")
		return
	}
	comma := false
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.IsExported() {
			if comma {
				b.WriteString("- ")
			}
			fmt.Fprintf(b, f.Name)
			fv := rv.Field(i)
			b.WriteRune(':')
			if fv.CanInterface() {
				printOn(fv.Interface(), b)
			}
			b.WriteString("\n")
			comma = true
		}
	}
}
