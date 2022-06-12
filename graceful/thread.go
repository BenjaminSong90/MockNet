package graceful

import "sync"

type routineGroup struct {
	waitGroup sync.WaitGroup
}


func createRoutineGroup() *routineGroup {
	return new(routineGroup)
}

func (g *routineGroup) Run(fun func())  {
	g.waitGroup.Add(1)
	go func() {
		defer g.waitGroup.Done()
		fun()
	}()
}

func (g *routineGroup) Wait() {
	g.waitGroup.Wait()
}
