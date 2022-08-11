package main

import (
	"context"
	"mocknet/fwatcher"
	"mocknet/graceful"
	"mocknet/logger"
	"mocknet/server"
	"mocknet/setting"
	"mocknet/utils"
)

func main() {
	logger.PlantTree(&logger.DebugTree{})

	utils.CheckModuleOrCreate()
	setting.LoadProjectConfig()

	manager := graceful.InitManager(context.Background())

	manager.AddRunningJob(func(ctx context.Context) error {
		server.StartServer(ctx)
		return nil
	})

	manager.AddRunningJob(func(ctx context.Context) error {
		localApiInfoPth := setting.GetLocalApiInfoPath()
		if setting.IsFileWatcherOpen() && len(localApiInfoPth) != 0 {
			fwatcher.InitLimit()
			fwatcher.GetFileWatcher().Watch(localApiInfoPth)
		}
		return nil
	})

	manager.AddShutdownJob(func() error {
		fwatcher.GetFileWatcher().Close()
		return nil
	})

	<-manager.Done()

}
