package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emicklei/simone/api"
	"github.com/peterh/liner"
)

type actionCommander struct {
	runner runnable
}

func newActionCommander(r runnable) *actionCommander {
	return &actionCommander{
		runner: r,
	}
}

const hist = ".simone"

func (a *actionCommander) Loop() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)

	if f, err := os.Open(hist); err == nil {
		line.ReadHistory(f)
		f.Close()
	}
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
				output(Print(res.RawData), true)
				continue
			}
			if entry == ":p" {
				res := a.RunString("_plugins()")
				output(Print(res.RawData), true)
				continue
			}
			if entry == ":d" {
				a.RunString("_toggledebug()")
				continue
			}
			if entry == ":h" {
				res := a.RunString("_showhelp()")
				output(Print(res.RawData), true)
				continue
			}
			if entry == ":w" {
				go func() { // blocks on exit UI
					if err := openSimonUI(); err != nil {
						output(err.Error(), false)
					}
				}()
				continue
			}
			if strings.HasPrefix(entry, ":l") {
				target := ""
				if arg := entry[2:]; len(arg) > 0 {
					target = strings.Trim(arg, " ")
				}
				if target == "" {
					input, err := line.Prompt("  plugin:")
					if err != nil {
						output(Print(err), false)
						continue
					}
					target = input
				}
				username, err := line.Prompt(fmt.Sprintf("  [%s] user:", target))
				if err != nil {
					output(Print(err), false)
					continue
				}
				password, err := line.PasswordPrompt(fmt.Sprintf("  [%s] password:", target))
				if err != nil {
					output(Print(err), false)
					continue
				}
				res := a.RunString(fmt.Sprintf("_login(%s,%q,%q)", target, username, password))
				if res.Error != "" {
					output(Print(res.Error), false)
				}
				continue
			}
			output("unknown command "+entry, false)
			continue
		}
		// ? = what can this value do
		if strings.HasSuffix(entry, "?") {
			src := fmt.Sprintf("_methods(%s)", entry[0:len(entry)-1])
			res := a.RunString(src)
			output(Print(res.RawData), true)
			continue
		}
		// ! = open browser on object
		if strings.HasSuffix(entry, "!") {
			src := fmt.Sprintf("_browse(%s)", entry[0:len(entry)-1])
			res := a.RunString(src)
			if res.Error != "" {
				output(Print(res.RawData), true)
			}
			continue
		}
		line.AppendHistory(entry)
		res := a.RunString(entry)
		if res.Error != "" {
			output(res.Error, false)
		} else {
			output(Print(res.RawData), true)
		}
	}
exit:
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

func (a *actionCommander) RunString(entry string) api.EvalResult {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\033[1;31mdo not panic:%v\033[0m\n", r)
		}
	}()
	return a.runner.RunString(entry)
}
