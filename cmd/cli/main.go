package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/emicklei/simone"
	"github.com/peterh/liner"
)

func main() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	client := new(http.Client)

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
				continue
			}
		}
		line.AppendHistory(entry)
		v, err := RunString(client, entry)
		if err != nil {
			fmt.Printf("\043[1;33m%v\043[0m\n", err)
		} else {
			fmt.Printf("\033[1;33m%v\033[0m\n", v)
		}
	}
exit:
}
func RunString(client *http.Client, entry string) (any, error) {
	body := bytes.NewBufferString(entry)
	req, err := http.NewRequest(http.MethodPost, "http://localhost:9119/v1/statements", body)
	if err != nil {
		return nil, err
	}
	ap := simone.ActionParams{
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
