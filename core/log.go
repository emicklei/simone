package core

import "fmt"

func loginfo(s string) {
	fmt.Printf("\033[1;36m%s\033[0m\n", s)
}
