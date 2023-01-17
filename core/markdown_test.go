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
	s = PrintMarkdown([]int{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 126, 27, 28, 29, 30})
	if got, want := flatten(s), "0. -1!1. !2. 1!3. 2!4. 3!5. 4!6. 5!7. 6!8. 7!9. 8!10. 9!11. 10!12. 11!13. 12!14. 13!15. 14!16. 15!17. 126!18. 27!19. 28!21. (2 more)!"; got != want {
		t.Errorf("got [%v:%T] want [%v:%T]", got, got, want, want)
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
