package simone

import (
	"fmt"
	"reflect"
	"strings"
)

func PrintStruct(v any) string {
	if v == nil {
		return "null"
	}
	b := new(strings.Builder)
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
	return b.String()
}
