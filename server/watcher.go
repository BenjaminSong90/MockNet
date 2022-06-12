package server

import (
	"context"
	"log"
	"mock_net/setting"
	"mock_net/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/howeyc/fsnotify"
)

func watchFolder(path string, ctx context.Context) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	go func(c context.Context) {
		for {
			select {
			case ev := <-watcher.Event:
				if isWatchedFile(ev.Name) {
					utils.DebugLogger("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				utils.DebugLogger("error: %s", err)
			case <-c.Done():
				return
			}
		}
	}(ctx)

	utils.DebugLogger("Watching %s", path)
	err = watcher.Watch(path)

	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()
	<- ctx.Done()
}

func Watch(ctx context.Context) {
	paths := setting.GetLocalApiInfoPath()
	for _, root := range paths {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
					return filepath.SkipDir
				}

				if isIgnoredFolder(path) {
					utils.DebugLogger("Ignoring %s", path)
					return filepath.SkipDir
				}

				go watchFolder(path, ctx)
			}

			return err
		})
	}

}

func isWatchedFile(path string) bool {
	ext := filepath.Ext(path)


	for _, e := range strings.Split(setting.GetFileWatcherValidExt(), ",") {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func isIgnoredFolder(path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	if len(setting.GetFileWatcherIgnoredFolder()) == 0{
		return false
	}

	for _, e := range strings.Split(setting.GetFileWatcherIgnoredFolder(), ",") {
		if strings.TrimSpace(e) == paths[0] {
			return true
		}
	}
	return false
}
