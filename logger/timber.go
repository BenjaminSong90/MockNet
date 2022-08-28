package logger

import (
	"fmt"
	"sync"
)

type Tree interface {
	V(format string)

	D(format string)

	I(format string)

	W(format string)

	E(format string)
}

type timber struct {
	forest []Tree
	sync.RWMutex
}

var _ Tree = &timber{}

var treeOfSouls = timber{}

func PlantTree(tree Tree) {
	treeOfSouls.Lock()
	defer treeOfSouls.Unlock()
	if tree != nil {
		treeOfSouls.forest = append(treeOfSouls.forest, tree)
	}
}

func UnRoot(tree Tree) {
	treeOfSouls.Lock()
	defer treeOfSouls.Unlock()
	for i := 0; i < len(treeOfSouls.forest); i++ {
		if tree == treeOfSouls.forest[i] {
			treeOfSouls.forest = append(treeOfSouls.forest[:i], treeOfSouls.forest[i+1:]...)
			break
		}
	}
}

func UnRootAll() {
	treeOfSouls.Lock()
	defer treeOfSouls.Unlock()
	treeOfSouls.forest = treeOfSouls.forest[:0]
}

func Size() int {
	treeOfSouls.Lock()
	defer treeOfSouls.Unlock()
	return len(treeOfSouls.forest)
}

func (t *timber) V(format string) {
	for _, tree := range t.forest {
		tree.V(format)
	}
}

func (t *timber) D(format string) {
	for _, tree := range t.forest {
		tree.D(format)
	}
}

func (t *timber) I(format string) {
	for _, tree := range t.forest {
		tree.I(format)
	}
}

func (t *timber) W(format string) {
	for _, tree := range t.forest {
		tree.W(format)
	}
}

func (t *timber) E(format string) {
	for _, tree := range t.forest {
		tree.E(format)
	}
}

type DebugTree struct{}

var _ Tree = &DebugTree{}

func (t *DebugTree) V(format string) {
	vLogger(format)
}

func (t *DebugTree) D(format string) {
	dLogger(format)
}

func (t *DebugTree) I(format string) {
	iLogger(format)
}

func (t *DebugTree) W(format string) {
	wLogger(format)
}

func (t *DebugTree) E(format string) {
	eLogger(format)
}

func V(format string, args ...interface{}) {
	treeOfSouls.V(fmt.Sprintf(format, args...))
}

func D(format string, args ...interface{}) {
	treeOfSouls.D(fmt.Sprintf(format, args...))
}

func I(format string, args ...interface{}) {
	treeOfSouls.I(fmt.Sprintf(format, args...))
}

func W(format string, args ...interface{}) {
	treeOfSouls.W(fmt.Sprintf(format, args...))
}

func E(format string, args ...interface{}) {
	treeOfSouls.E(fmt.Sprintf(format, args...))
}
