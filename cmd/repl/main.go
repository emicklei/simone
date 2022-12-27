package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

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
				v, _ := RunString(client, "_variables()")
				fmt.Printf("\033[1;33m%v\033[0m\n", v)
				continue
				continue
			}
			if entry == ":p" {
				v, _ := RunString(client, "_plugins()")
				fmt.Printf("\033[1;33m%v\033[0m\n", v)
				continue
			}
		}
		line.AppendHistory(entry)
		v, err := RunString(client, entry)
		if err != nil {
			fmt.Printf("\033[1;31m%v\033[0m\n", err)
		} else {
			fmt.Printf("\033[1;33m%v\033[0m\n", v)
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
func RunString(client *http.Client, entry string) (any, error) {
	body := bytes.NewBufferString(entry)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9119/v1", body)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return string(result), nil
}
