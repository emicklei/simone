package simone

import (
	"strings"
	"testing"
)

func TestPrintNil(t *testing.T) {
	s := Print(nil)
	if got, want := s, "null"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	s = Print(12)
	if got, want := s, "12"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	s = Print("string")
	if got, want := s, "'string'"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	type Range struct {
		Low       int
		High      int
		Inclusive *bool
	}
	w := true
	s = Print(&Range{-1, 1, &w})
	if got, want := s, "{'Low':-1, 'High':1, 'Inclusive':true}"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
	RegisterPrinter(new(Range), RangePrinter)
	s = Print(&Range{-1, 1, &w})
	if got, want := s, "rangeprinted"; got != want {
		t.Errorf("got [%v]:%T want [%v]:%T", got, got, want, want)
	}
}

func RangePrinter(v any, b *strings.Builder) {
	b.WriteString("rangeprinted")
}
