package server

import (
	"context"
	"sync"
)

var (
	startChannel chan string
	stopChannel  chan bool
)

func init() {
	startChannel = make(chan string, 1000)
	stopChannel = make(chan bool)
}

func StartServer(ctx context.Context)  {
	var wg sync.WaitGroup
	go func(c context.Context) {
		wg.Add(1)
		defer wg.Done()
		for{
			var container  = Container{}
			container.Start(New())
			select {
				case reStart :=<- stopChannel:
					container.CloseWithWait()
					if !reStart {
						return
					}
				case <- c.Done():
					container.CloseWithWait()
					return
			}

		}
	}(ctx)
	<- ctx.Done()
	wg.Wait()
}