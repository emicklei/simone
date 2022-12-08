package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/emicklei/simone/api"
	"github.com/peterh/liner"
	"google.golang.org/grpc"
)

func main() {
	line := liner.NewLiner()
	line.SetCtrlCAborts(true)
	conn, err := grpc.Dial("localhost:9191", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	eval := api.NewEvaluationServiceClient(conn)

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
		v, err := RunString(eval, entry)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("\033[1;33m%v\033[0m\n", v)
		}
	}
exit:
}
func RunString(client api.EvaluationServiceClient, entry string) (any, error) {
	req := new(api.EvalRequest)
	req.Source = entry
	repl, err := client.Eval(context.Background(), req)
	if err != nil {
		return nil, err
	}
	res := repl.Result
	return res, nil
}
