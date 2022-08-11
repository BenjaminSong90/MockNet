package fwatcher

import (
	"context"
	"mocknet/logger"
	"mocknet/setting"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/howeyc/fsnotify"
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
			case ev, ok := <-watcher.Event:
				if ok {
					logger.D("sending event %s", ev.String())
					if fw.isWatchedFile(ev.Name) {
						FileChangeChannel <- ev.String()
					}
				} else {
					logger.D("sending empty event")
				}
			case err, ok := <-watcher.Error:
				if ok {
					logger.D("watcher error: %s", err)
				}
			case <-c.Done():
				goto OutLoop
			}
		}
	OutLoop:
	}(fw.ctx)

	logger.D("Watching %s", path)
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
					logger.D("Ignoring %s", path)
					return filepath.SkipDir
				}

				w, err := fw.watchFolder(path)
				if err != nil {
					logger.E("%s error: %s", path, err)
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
	logger.D("Start close file watcher")
	fileWatcher.ctxCancel()
	var errs []error
	for _, w := range fw.watchers {
		err := w.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	logger.D("End close file watcher")
	close(FileChangeChannel)
	return errs
}

func GetFileWatcher() *FileWatcher {
	startOnce.Do(func() {
		fileWatcher = &FileWatcher{
			lock:     &sync.RWMutex{},
			watchers: make([]*fsnotify.Watcher, 0),
		}
		ctx, cancel := context.WithCancel(context.Background())
		fileWatcher.ctx = ctx
		fileWatcher.ctxCancel = cancel
	})

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
