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
	handler *ActionHandler
}

func NewActionCommander(h *ActionHandler) *ActionCommander {
	return &ActionCommander{
		handler: h,
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

	fmt.Printf("\033[1;32m%s\033[0m\n", ":q (quit) :p (plugins) :v (variables)")

	defer line.Close()
	for {
		entry, err := line.Prompt("⨁ ")
		if err != nil {
			log.Println(err)
			goto exit
		}
		if strings.HasPrefix(entry, ":") {
			// special case
			if entry == ":q" || entry == ":Q" {
				goto exit
			}
			if entry == ":v" {
				res := a.RunString("_variables()")
				fmt.Printf("\033[1;33m%v\033[0m\n", res.Data)
				continue
			}
			if entry == ":p" {
				res := a.RunString("_plugins()")
				fmt.Printf("\033[1;33m%v\033[0m\n", res.Data)
				continue
			}
		}
		line.AppendHistory(entry)
		res := a.RunString(entry)
		if res.Error != "" {
			fmt.Printf("\033[1;31m%v\033[0m\n", res.Error)
		} else {
			fmt.Printf("\033[1;33m%v\033[0m\n", res.Data)
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

func (a *ActionCommander) RunString(entry string) api.EvalResult {
	// TODO dedup code
	res := api.EvalResult{}
	result, err := a.handler.vm.RunString(entry)
	if err != nil {
		// syntax error
		res.Error = err.Error()
	} else {
		val := result.Export()
		if err, ok := val.(error); ok {
			// evaluation error
			res.Error = err.Error()
		} else {
			// no error
			res.Data = Print(val)
			res.Datatype = fmt.Sprintf("%T", val)
		}
	}
	return res
}
