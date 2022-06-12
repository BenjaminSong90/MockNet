package main

import (
	"context"
	"mock_net/graceful"
	"mock_net/server"
	"mock_net/setting"
)

func main() {

	setting.LoadProjectConfig()
	ctx, _ := context.WithCancel(context.Background())

	manager := graceful.InitManager(ctx)

	manager.AddRunningJob(func(ctx context.Context) error {
		server.StartServer(ctx)
		return nil
	})

	manager.AddRunningJob(func(ctx context.Context) error {
		if setting.IsFileWatcherOpen() {//TODO 需要优化关闭方式
			server.InitLimit()
			server.Watch(ctx)
		}
		return nil
	})


	<-manager.Done()


}
