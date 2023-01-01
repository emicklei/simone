package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emicklei/simone/api"
	"github.com/peterh/liner"
)

type ActionCommander struct {
	runner Runnable
}

func NewActionCommander(r Runnable) *ActionCommander {
	return &ActionCommander{
		runner: r,
	}
}

const hist = ".simone"

func (a *ActionCommander) Loop() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	if f, err := os.Open(hist); err == nil {
		line.ReadHistory(f)
		f.Close()
	}

	fmt.Printf("\033[1;32m%s\033[0m\n", ":q (quit) :h (help) :v (variables) :p (plugins) :d (verbose)")

	defer line.Close()
	for {
		entry, err := line.Prompt("⨁ ")
		if err != nil {
			log.Println(err)
			goto exit
		}
		if strings.HasPrefix(entry, ":") {
			// special cases
			if entry == ":q" || entry == ":Q" {
				goto exit
			}
			if entry == ":v" {
				res := a.RunString("_variables()")
				output(res.Data, true)
				continue
			}
			if entry == ":p" {
				res := a.RunString("_plugins()")
				output(res.Data, true)
				continue
			}
			if entry == ":d" {
				a.RunString("_toggledebug()")
				continue
			}
			if entry == ":h" {
				res := a.RunString("_showhelp()")
				output(res.Data, true)
				continue
			}
		}
		// ? = what can this value do
		if strings.HasSuffix(entry, "?") {
			src := fmt.Sprintf("_methods(%s)", entry[0:len(entry)-1])
			res := a.RunString(src)
			output(res.Data, true)
			continue
		}
		line.AppendHistory(entry)
		res := a.RunString(entry)
		if res.Error != "" {
			output(res.Error, false)
		} else {
			output(res.Data, true)
		}
	}
exit:
	fmt.Printf("\033[1;32m%s\033[0m\n", "simone says bye!")
	if f, err := os.Create(hist); err != nil {
		log.Print("Error writing history file: ", err)
	} else {
		line.WriteHistory(f)
		f.Close()
	}
}

func output(v any, ok bool) {
	if !ok {
		fmt.Printf("\033[1;31m%v\033[0m\n", v)
	} else {
		fmt.Printf("\033[1;33m%v\033[0m\n", v)
	}
}

func (a *ActionCommander) RunString(entry string) api.EvalResult {
	return a.runner.RunString(entry)
}
