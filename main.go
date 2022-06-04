package main

import (
	"context"
	"mock_net/server"
	"mock_net/utils"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	go utils.ListenBreak(cancel)
	server.StartServer(ctx)
}
