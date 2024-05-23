package main

import (
	"context"
	"os"
	"syscall"

	"github.com/oklog/run"
	"github.com/peizhong/codeplay/service/lambda"
)

func main() {
	r := &run.Group{}
	r.Add(run.SignalHandler(context.TODO(), os.Interrupt, syscall.SIGTERM))
	lambda.RegisterRuntime(r)
	r.Run()
}
