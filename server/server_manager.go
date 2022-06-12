package server

import (
	"context"
	"mocknet/fwatcher"
	"mocknet/setting"
	"mocknet/utils"
	"strings"
	"time"
)

var (
	stopChannel  chan bool
)

func init() {
	stopChannel = make(chan bool)
}

func StartServer(ctx context.Context)  {
	for{
		var container  = Container{}
		container.Start(New())
		select {
		case reStart :=<- stopChannel:
			container.CloseWithWait()
			if !reStart {
				return
			}
		case eventName :=<- fwatcher.FileChangeChannel:

			time.Sleep(1000 * time.Millisecond)

			flushEvents()

			if shouldReload(eventName){
				container.CloseWithWait()
			}

		case <- ctx.Done():
			container.CloseWithWait()
			return
		}

	}
}


func flushEvents() {
	for {
		select {
		case eventName := <- fwatcher.FileChangeChannel:
			utils.DebugLogger("receiving event %s", eventName)
		default:
			return
		}
	}
}

func shouldReload(eventName string) bool {
	for _, e := range strings.Split(setting.GetFileWatcherNoReloadExt(), ",") {
		e = strings.TrimSpace(e)
		fileName := strings.Replace(strings.Split(eventName, ":")[0], `"`, "", -1)
		if strings.HasSuffix(fileName, e) {
			return false
		}
	}

	return true
}