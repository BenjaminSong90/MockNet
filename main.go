package main

import (
	"context"
	"mock_net/server"
	"mock_net/setting"
	"mock_net/utils"
)

func main() {

	setting.LoadProjectConfig()
	ctx, cancel := context.WithCancel(context.Background())
	if setting.IsFileWatcherOpen() {
		server.InitLimit()
		server.Watch(ctx)
	}

	go utils.ListenBreak(cancel)
	server.StartServer(ctx)
}
