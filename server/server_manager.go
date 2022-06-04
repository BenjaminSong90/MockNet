package server

import (
	"context"
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
	for{
		var container  = Container{}
		container.Start(New())
		select {
		case reStart :=<- stopChannel:
			container.CloseWithWait()
			if !reStart {
				return
			}
		case <- ctx.Done():
			container.CloseWithWait()
			return
		}

	}
}