package core

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/emicklei/simone/api"
)

func buildInspectResult(res api.EvalResult) api.InspectResult {
	if res.RawData == nil {
		// TODO how to inspect null
		m := map[string]any{"_": "null"}
		return api.InspectResult{Datatype: "?", Object: m}
	}
	if res.Error != "" {
		return api.InspectResult{
			Error:    res.Error,
			Datatype: res.Datatype,
		}
	}
	ires := api.InspectResult{Datatype: fmt.Sprintf("%T", res.RawData)}
	switch v := res.RawData.(type) {
	case []any:
		m := map[string]any{}
		for i, each := range v {
			m[strconv.Itoa(i)] = each
		}
		ires.Object = m
	case map[string]any:
		ires.Object = v
	default:
		// check underlying type is map[string]any
		rt := reflect.TypeOf(res.RawData)
		rv := reflect.ValueOf(res.RawData)
		if rt.Kind() == reflect.Pointer {
			rt = rt.Elem()
			rv = rv.Elem()
		}
		if rt.Kind() == reflect.Map {
			if rt.Key().Kind() == reflect.String {
				if rt.Elem().Kind() == reflect.Interface {
					m := map[string]any{}
					rv := reflect.ValueOf(res.RawData)
					iter := rv.MapRange()
					for iter.Next() {
						m[iter.Key().String()] = valueOrPrintstring(iter.Value().Interface())
					}
					ires.Object = m
					return ires
				}
			}
		}
		if rt.Kind() == reflect.Struct {
			m := map[string]any{}
			for i := 0; i < rt.NumField(); i++ {
				f := rt.Field(i)
				if f.IsExported() {
					fv := rv.Field(i)
					if fv.CanInterface() {
						m[f.Name] = valueOrPrintstring(rv.Field(i).Interface())
					}
				}
			}
			ires.Object = m
			return ires
		}
		if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
			m := map[string]any{}
			for i := 0; i < rv.Len(); i++ {
				ev := rv.Index(i)
				if ev.CanInterface() {
					m[strconv.Itoa(i)] = valueOrPrintstring(ev.Interface())
				}
			}
			ires.Object = m
			return ires
		}
		ires.Object = map[string]any{
			"_": res.RawData,
		}
	}
	return ires
}

// valueOrPrintstring checks registered printers
func valueOrPrintstring(v any) any {
	rt := reflect.TypeOf(v)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}
	if pf, ok := printer.registry[rt]; ok {
		b := new(strings.Builder)
		pf(v, b)
		return b.String()
	}
	return v
}
