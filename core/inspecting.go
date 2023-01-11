package core

import (
	"strconv"

	"github.com/emicklei/simone/api"
)

func buildInspectResult(res api.EvalResult) api.InspectResult {
	if res.Error != "" {
		return api.InspectResult{
			Error:    res.Error,
			Datatype: res.Datatype,
		}
	}
	ires := api.InspectResult{}
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
		ires.Object = map[string]any{
			"_": res.RawData,
		}
	}
	return ires
}
