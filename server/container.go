package server

import (
	"context"
	"sync"
)

type ContainerRunnable interface {
	Run(ctx context.Context) error
}

type Container struct {
	once       sync.Once
	ctx        context.Context
	cancelFunc context.CancelFunc
	waitGroup  sync.WaitGroup
	runnable   ContainerRunnable
}

func (ctn *Container) Start(runnable ContainerRunnable) {
	ctn.once.Do(func() {
		ctn.runnable = runnable
		ctn.ctx, ctn.cancelFunc = context.WithCancel(context.Background())
		go func() {
			defer ctn.waitGroup.Done()
			ctn.waitGroup.Add(1)
			err := runnable.Run(ctn.ctx)
			if err != nil {
				panic(err)
			}

		}()
	})
}

func (ctn *Container) Wait() {
	ctn.waitGroup.Wait()
}

func (ctn *Container) Close() {
	ctn.cancelFunc()
}

func (ctn *Container) CloseWithWait() {
	ctn.Close()
	ctn.Wait()
}
