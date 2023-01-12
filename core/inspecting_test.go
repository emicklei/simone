package core

import (
	"testing"

	"github.com/emicklei/simone/api"
)

func TestBuildInspect_CustomMap(t *testing.T) {
	type mymap map[string]any
	mm := mymap{"test": 1}
	res := api.EvalResult{RawData: mm}
	ires := buildInspectResult(res)
	if got, want := ires.Object["test"], 1; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
func TestBuildInspect_Int(t *testing.T) {
	res := api.EvalResult{RawData: 42}
	ires := buildInspectResult(res)
	if got, want := ires.Object["_"], 42; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestBuildInspect_Struct(t *testing.T) {
	type mystruct struct{ Field string }
	ms := mystruct{Field: "field"}
	res := api.EvalResult{RawData: ms}
	ires := buildInspectResult(res)
	if got, want := ires.Object["Field"], "field"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	ms2 := &mystruct{Field: "field"}
	res2 := api.EvalResult{RawData: ms2}
	ires2 := buildInspectResult(res2)
	if got, want := ires2.Object["Field"], "field"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func TestBuildInspect_Array(t *testing.T) {
	a := []int{1, 2}
	res := api.EvalResult{RawData: a}
	ires := buildInspectResult(res)
	if got, want := ires.Datatype, "[]int"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := ires.Object["1"], 2; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	a2 := [2]int{1, 2}
	res2 := api.EvalResult{RawData: a2}
	ires2 := buildInspectResult(res2)
	if got, want := ires2.Datatype, "[2]int"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	if got, want := ires2.Object["1"], 2; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}
