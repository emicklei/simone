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
	insp := api.NewInspectServiceClient(conn)

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
		v, err := RunString(eval, insp, entry)
		if err != nil {
			fmt.Printf("\043[1;33m%v\043[0m\n", err)
		} else {
			fmt.Printf("\033[1;33m%v\033[0m\n", v)
		}
	}
exit:
}
func RunString(eval api.EvaluationServiceClient,
	insp api.InspectServiceClient,
	entry string) (any, error) {
	req := new(api.EvalRequest)
	req.Source = entry
	repl, err := eval.Eval(context.Background(), req)
	if err != nil {
		return nil, err
	}
	req2 := new(api.PrintStringRequest)
	req2.Id = repl.RemoteObjectId
	repl2, err := insp.PrintString(context.Background(), req2)
	if err != nil {
		return nil, err
	}
	return repl2.Content, nil
}
