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

var tree_of_souls = timber{}

func PlantTree(tree Tree) {
	tree_of_souls.Lock()
	defer tree_of_souls.Unlock()
	if tree != nil {
		tree_of_souls.forest = append(tree_of_souls.forest, tree)
	}
}

func UnRoot(tree Tree) {
	tree_of_souls.Lock()
	defer tree_of_souls.Unlock()
	for i := 0; i < len(tree_of_souls.forest); i++ {
		if tree == tree_of_souls.forest[i] {
			tree_of_souls.forest = append(tree_of_souls.forest[:i], tree_of_souls.forest[i+1:]...)
			break
		}
	}
}

func UnRootAll() {
	tree_of_souls.Lock()
	defer tree_of_souls.Unlock()
	tree_of_souls.forest = tree_of_souls.forest[:0]
}

func Size() int {
	tree_of_souls.Lock()
	defer tree_of_souls.Unlock()
	return len(tree_of_souls.forest)
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
	tree_of_souls.V(fmt.Sprintf(format, args...))
}

func D(format string, args ...interface{}) {
	tree_of_souls.D(fmt.Sprintf(format, args...))
}

func I(format string, args ...interface{}) {
	tree_of_souls.I(fmt.Sprintf(format, args...))
}

func W(format string, args ...interface{}) {
	tree_of_souls.W(fmt.Sprintf(format, args...))
}

func E(format string, args ...interface{}) {
	tree_of_souls.E(fmt.Sprintf(format, args...))
}
