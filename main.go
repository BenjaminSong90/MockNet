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
	fw := server.InitFileWatcher(ctx)

	manager.AddRunningJob(func(ctx context.Context) error {
		server.StartServer(ctx)
		return nil
	})

	manager.AddRunningJob(func(ctx context.Context) error {
		localApiInfoPth := setting.GetLocalApiInfoPath()
		if setting.IsFileWatcherOpen() && len(localApiInfoPth) != 0{
			server.InitLimit()
			fw.Watch(localApiInfoPth)
		}
		return nil
	})

	manager.AddShutdownJob(func() error {
		fw.Close()
		return nil
	})


	<-manager.Done()


}
