package server

import (
	"context"
	"mock_net/setting"
	"mock_net/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/howeyc/fsnotify"
)

type FileWatcher struct {
	lock     *sync.RWMutex
	watchers []*fsnotify.Watcher
	ctx      context.Context
}

var fileWatcher *FileWatcher

var startOnce = sync.Once{}

func (fw *FileWatcher) watchFolder(path string) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func(c context.Context) {
		for {
			select {
			case ev := <-watcher.Event:
				if fw.isWatchedFile(ev.Name) {
					utils.DebugLogger("sending event %s", ev)
					startChannel <- ev.String()
				}
			case err := <-watcher.Error:
				utils.DebugLogger("error: %s", err)
			case <-c.Done():
				goto OutLoop
			}
		}
	OutLoop:
	}(fw.ctx)

	utils.DebugLogger("Watching %s", path)
	err = watcher.Watch(path)

	if err == nil {
		return watcher, nil
	}
	return nil, err
}

func (fw *FileWatcher) Watch(paths []string) {
	if len(paths) == 0 {
		panic("Api file paths is not empty")
	}
	for _, root := range paths {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				if len(path) > 1 && strings.HasPrefix(filepath.Base(path), ".") {
					return filepath.SkipDir
				}

				if fw.isIgnoredFolder(path) {
					utils.DebugLogger("Ignoring %s", path)
					return filepath.SkipDir
				}

				w, err := fw.watchFolder(path)
				if err != nil {
					utils.ErrorLogger("%s error: %s", path, err)
				} else {
					fw.lock.Lock()
					fw.watchers = append(fw.watchers, w)
					fw.lock.Unlock()
				}
			}

			return err
		})
	}
}

func (fw *FileWatcher) Close() {
	var errs []error
	for _, w := range fw.watchers {
		err := w.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	utils.DebugLogger("Close file watcher")
}

func InitFileWatcher(ctx context.Context) *FileWatcher {
	startOnce.Do(func() {
		fileWatcher = &FileWatcher{
			lock:     &sync.RWMutex{},
			ctx:      ctx,
			watchers: make([]*fsnotify.Watcher, 0),
		}
	})

	return fileWatcher
}

func GetFileWatcher() *FileWatcher {
	if fileWatcher == nil {
		panic("File watcher is not init")
	}

	return fileWatcher
}

func (fw *FileWatcher) isWatchedFile(path string) bool {
	ext := filepath.Ext(path)

	for _, e := range strings.Split(setting.GetFileWatcherValidExt(), ",") {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func (fw *FileWatcher) isIgnoredFolder(path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	if len(setting.GetFileWatcherIgnoredFolder()) == 0 {
		return false
	}

	for _, e := range strings.Split(setting.GetFileWatcherIgnoredFolder(), ",") {
		if strings.TrimSpace(e) == paths[0] {
			return true
		}
	}
	return false
}
