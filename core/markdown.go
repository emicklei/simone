package core

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"github.com/emicklei/simone/api"
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
	if rt.Kind() == reflect.Map {
		printMapMarkdown(b, rt, rv)
		return
	}
	// fallback to standard JSON encoder
	if api.Debug {
		log.Println("markdown fallback JSON")
	}
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
	b.WriteString(rt.PkgPath())
	b.WriteRune('/')
	b.WriteString(rt.Name())
	b.WriteString("\n\n")
	if !rv.IsValid() {
		b.WriteString("null")
		return
	}
	kvs := []kv{}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if f.IsExported() {
			fv := rv.Field(i)
			if fv.CanInterface() {
				kvs = append(kvs, kv{k: f.Name, v: fv.Interface()})
			}
		}
	}
	printKVsOn(kvs, b)
}

func printMapMarkdown(b *strings.Builder, rt reflect.Type, rv reflect.Value) {
	r := rv.MapRange()
	kvs := []kv{}
	for r.Next() {
		k := fmt.Sprintf("%v", r.Key())
		v := r.Value()
		if v.CanInterface() {
			kvs = append(kvs, kv{k: k, v: v.Interface()})
		}
	}
	printKVsOn(kvs, b)
}

type kv struct {
	k string
	v any
}

func printKVsOn(kvs []kv, b *strings.Builder) {
	sort.Slice(kvs, func(i, j int) bool { return kvs[i].k < kvs[j].k })
	for _, each := range kvs {
		fmt.Fprintf(b, "- %s: %v\n", each.k, each.v)
	}
}
