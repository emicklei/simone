package core

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/emicklei/simone/api"
)

type CanBeObject interface {
	ToObject() map[string]any
}

func buildInspectResult(res api.EvalResult) api.InspectResult {
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
		if rt.Kind() == reflect.Map {
			if rt.Key().Kind() == reflect.String {
				if rt.Elem().Kind() == reflect.Interface {
					m := map[string]any{}
					rv := reflect.ValueOf(res.RawData)
					iter := rv.MapRange()
					for iter.Next() {
						m[iter.Key().String()] = iter.Value().Interface()
					}
					ires.Object = m
					return ires
				}
			}
		}
		ires.Object = map[string]any{
			"_": res.RawData,
		}
	}
	return ires
}
