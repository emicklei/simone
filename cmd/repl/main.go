package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/emicklei/simone/api"
	"github.com/emicklei/simone/core"
	"github.com/peterh/liner"
)

const hist = ".simone"

func main() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	client := new(http.Client)

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
			// special case
			if entry == ":q" || entry == ":Q" {
				goto exit
			}
			if entry == ":v" {
				res := RunString(client, "_variables()")
				fmt.Printf("\033[1;33m%v\033[0m\n", res.Data)
				continue
			}
			if entry == ":p" {
				res := RunString(client, "_plugins()")
				fmt.Printf("\033[1;33m%v\033[0m\n", res.Data)
				continue
			}
		}
		line.AppendHistory(entry)
		res := RunString(client, entry)
		if res.Error != "" {
			fmt.Printf("\033[1;31m%v\033[0m\n", res.Error)
		} else {
			fmt.Printf("\033[1;33m%v\033[0m\n", res.Data)
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
func RunString(client *http.Client, entry string) api.EvalResult {
	body := bytes.NewBufferString(entry)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9119/v1", body)
	res := api.EvalResult{}
	if err != nil {
		res.Error = err.Error()
		return res
	}
	ap := core.ActionParams{
		Debug:  true,
		Line:   "1",
		File:   "cli",
		Action: "eval",
	}
	ap.Inject(req.URL)
	resp, err := client.Do(req)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Error = err.Error()
		return res
	}
	if err := json.Unmarshal(result, &res); err != nil {
		res.Error = err.Error()
	}
	return res
}
