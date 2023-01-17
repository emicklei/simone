package core

import (
	"strings"
	"testing"
)

func TestPrintMarkdown(t *testing.T) {
	s := PrintMarkdown(nil)
	if got, want := s, "null"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	s = PrintMarkdown(42)
	if got, want := s, "42"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
	s = PrintMarkdown([]int{-1, 0, 1})
	if got, want := flatten(s), "0. -1!1. 0!2. 1!"; got != want {
		t.Errorf("got [\n%v:%T] want [\n%v:%T]", got, got, want, want)
	}
	w := false
	r := &Range{-1, 1, &w}
	s = PrintMarkdown(r)
	if got, want := flatten(s), "github.com/emicklei/simone/core/Range!!- High: 1!- Inclusive: false!- Low: -1!"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
	}
}

func flatten(s string) string {
	return strings.ReplaceAll(s, "\n", "!")
}
