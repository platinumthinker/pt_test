package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/platinumthinker/pt_test/src/server/processor"
	"github.com/platinumthinker/pt_test/src/server/transport"
)

const (
	listenAddress = "127.0.0.1:8080"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	jsonProcessor := processor.NewJsonProcessor()
	jsonProcessor.AddOperation("sum", &processor.SumOperation{})
	jsonProcessor.AddOperation("mul", &processor.MulOperation{})
	listener := transport.NewWebSocketListener(listenAddress)
	listener.Handle("/ws", jsonProcessor)

	err := listener.Listen()
	if err != nil {
		fmt.Printf("error: transport problem: %v", err)
		return
	}

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		listener.Close(ctx)
		cancel()
	}()

	<-ctx.Done()
}
