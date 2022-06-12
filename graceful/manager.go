package graceful

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type (
	RunningJob func(context.Context) error
	ShutdownJob func() error
)

type Manager struct {
	lock              *sync.RWMutex
	shutdownCtx       context.Context
	shutdownCtxCancel context.CancelFunc
	waitCtx           context.Context
	waitCtxCancel     context.CancelFunc
	runningWaitGroup  *routineGroup
	errors            []error
	runAtShutdown     []ShutdownJob

}

var manager *Manager


func (g *Manager) handleSignals(ctx context.Context) {
	c := make(chan os.Signal, 1)

	signal.Notify(
		c,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)
	defer signal.Stop(c)

	for {
		select {
		case sig := <-c:
			switch sig {
			case syscall.SIGINT:
				g.doGracefulShutdown()
				return
			case syscall.SIGTERM:
				g.doGracefulShutdown()
				return
			case syscall.SIGTSTP:
			default:
			}
		case <-ctx.Done():
			g.doGracefulShutdown()
			return
		}
	}
}

func (g *Manager) start(ctx context.Context) {
	g.shutdownCtx, g.shutdownCtxCancel = context.WithCancel(ctx)
	g.waitCtx, g.waitCtxCancel = context.WithCancel(context.Background())

	go g.handleSignals(ctx)
}

func (g *Manager) doGracefulShutdown() {
	g.shutdownCtxCancel()
	for _, f := range g.runAtShutdown {
		func(run ShutdownJob) {
			g.runningWaitGroup.Run(func() {
				g.doShutdownJob(run)
			})
		}(f)
	}
	go func() {
		g.waitForJobs()
		g.lock.Lock()
		g.waitCtxCancel()
		g.lock.Unlock()
	}()
}

func (g *Manager) waitForJobs() {
	g.runningWaitGroup.Wait()
}

func (g *Manager) doShutdownJob(f ShutdownJob) {
	defer func() {
		if err := recover(); err != nil {
			msg := fmt.Errorf("panic in shutdown job: %v", err)
			g.lock.Lock()
			g.errors = append(g.errors, msg)
			g.lock.Unlock()
		}
	}()
	if err := f(); err != nil {
		g.lock.Lock()
		g.errors = append(g.errors, err)
		g.lock.Unlock()
	}
}

func (g *Manager) AddShutdownJob(f ShutdownJob) {
	g.lock.Lock()
	g.runAtShutdown = append(g.runAtShutdown, f)
	g.lock.Unlock()
}

func (g *Manager) AddRunningJob(f RunningJob) {
	g.runningWaitGroup.Run(func() {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Errorf("panic in running job: %v", err)
				g.lock.Lock()
				g.errors = append(g.errors, msg)
				g.lock.Unlock()
			}
		}()
		if err := f(g.shutdownCtx); err != nil {
			g.lock.Lock()
			g.errors = append(g.errors, err)
			g.lock.Unlock()
		}
	})
}


func (g *Manager) Done() <-chan struct{} {
	return g.waitCtx.Done()
}

func (g *Manager) ShutdownContext() context.Context {
	return g.shutdownCtx
}

var startOnce = sync.Once{}

func InitManager(ctx context.Context) *Manager {

	startOnce.Do(func() {
		manager =  &Manager{
			lock:             &sync.RWMutex{},
			errors:           make([]error, 0),
			runningWaitGroup: createRoutineGroup(),
		}
		manager.start(ctx)
	})

	return manager
}

func GetManager() *Manager {
	if manager == nil {
		panic("Manager is not init")
	}

	return manager
}