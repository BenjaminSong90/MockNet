package main

import (
	"context"
	"mocknet/fwatcher"
	"mocknet/graceful"
	"mocknet/server"
	"mocknet/setting"
	"mocknet/utils"
)

func main() {
	utils.CheckModuleOrCreate()
	setting.LoadProjectConfig()
	ctx, _ := context.WithCancel(context.Background())

	manager := graceful.InitManager(ctx)
	fw := fwatcher.InitFileWatcher()

	manager.AddRunningJob(func(ctx context.Context) error {
		server.StartServer(ctx)
		return nil
	})

	manager.AddRunningJob(func(ctx context.Context) error {
		localApiInfoPth := setting.GetLocalApiInfoPath()
		if setting.IsFileWatcherOpen() && len(localApiInfoPth) != 0{
			fwatcher.InitLimit()
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
