package logger

import "sync"

type Tree interface {
	V(format string, args ...interface{})

	D(format string, args ...interface{})

	I(format string, args ...interface{})

	W(format string, args ...interface{})

	E(format string, args ...interface{})
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

func (t *timber) V(format string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.V(format, args...)
	}
}

func (t *timber) D(format string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.D(format, args...)
	}
}

func (t *timber) I(format string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.I(format, args...)
	}
}

func (t *timber) W(format string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.W(format, args...)
	}
}

func (t *timber) E(format string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.E(format, args...)
	}
}

type DebugTree struct{}

var _ Tree = &DebugTree{}

func (t *DebugTree) V(format string, args ...interface{}) {
	vLogger(format, args)
}

func (t *DebugTree) D(format string, args ...interface{}) {
	dLogger(format, args)
}

func (t *DebugTree) I(format string, args ...interface{}) {
	iLogger(format, args)
}

func (t *DebugTree) W(format string, args ...interface{}) {
	wLogger(format, args)
}

func (t *DebugTree) E(format string, args ...interface{}) {
	eLogger(format, args)
}

func V(format string, args ...interface{}) {
	tree_of_souls.V(format, args...)
}

func D(format string, args ...interface{}) {
	tree_of_souls.D(format, args...)
}

func I(format string, args ...interface{}) {
	tree_of_souls.I(format, args...)
}

func W(format string, args ...interface{}) {
	tree_of_souls.W(format, args...)
}

func E(format string, args ...interface{}) {
	tree_of_souls.E(format, args...)
}
