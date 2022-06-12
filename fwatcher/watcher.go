package fwatcher

import (
	"context"
	"github.com/howeyc/fsnotify"
	"mocknet/setting"
	"mocknet/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	FileChangeChannel chan string
)

func init() {
	FileChangeChannel = make(chan string, 1000)
}

type FileWatcher struct {
	lock      *sync.RWMutex
	watchers  []*fsnotify.Watcher
	ctx       context.Context
	ctxCancel context.CancelFunc
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
				if ev != nil {
					utils.DebugLogger("sending event %s", ev.String())
					if fw.isWatchedFile(ev.Name) {
						FileChangeChannel <- ev.String()
					}
				} else {
					utils.DebugLogger("sending empty event")
				}
			case err := <-watcher.Error:
				if err != nil {
					utils.DebugLogger("watcher error: %s", err)
				}
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

func (fw *FileWatcher) Watch(paths []string) []error {
	if len(paths) == 0 {
		panic("Api file paths is not empty")
	}
	errs := make([]error, 0)
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
					errs = append(errs, err)
				} else {
					fw.lock.Lock()
					fw.watchers = append(fw.watchers, w)
					fw.lock.Unlock()
				}
			}

			return err
		})
	}
	return errs
}

func (fw *FileWatcher) Close() []error {
	utils.DebugLogger("Start close file watcher")
	fileWatcher.ctxCancel()
	var errs []error
	for _, w := range fw.watchers {
		err := w.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	utils.DebugLogger("End close file watcher")
	close(FileChangeChannel)
	return errs
}

func InitFileWatcher() *FileWatcher {
	startOnce.Do(func() {
		fileWatcher = &FileWatcher{
			lock:     &sync.RWMutex{},
			watchers: make([]*fsnotify.Watcher, 0),
		}
		ctx,cancel := context.WithCancel(context.Background())
		fileWatcher.ctx = ctx
		fileWatcher.ctxCancel = cancel
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
