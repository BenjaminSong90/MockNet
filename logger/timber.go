package logger

import "sync"

type Tree interface {
	V(err error, message string, args ...interface{})

	D(err error, message string, args ...interface{})

	I(err error, message string, args ...interface{})

	W(err error, message string, args ...interface{})

	E(err error, message string, args ...interface{})
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

func (t *timber) V(err error, message string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.V(err, message, args...)
	}
}

func (t *timber) D(err error, message string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.D(err, message, args...)
	}
}

func (t *timber) I(err error, message string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.I(err, message, args...)
	}
}

func (t *timber) W(err error, message string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.W(err, message, args...)
	}
}

func (t *timber) E(err error, message string, args ...interface{}) {
	for _, tree := range t.forest {
		tree.E(err, message, args...)
	}
}

type DebugTree struct{}

var _ Tree = &DebugTree{}

func (t *DebugTree) V(err error, message string, args ...interface{}) {
	vLogger(message, args)
}

func (t *DebugTree) D(err error, message string, args ...interface{}) {
	dLogger(message, args)
}

func (t *DebugTree) I(err error, message string, args ...interface{}) {
	iLogger(message, args)
}

func (t *DebugTree) W(err error, message string, args ...interface{}) {
	wLogger(message, args)
}

func (t *DebugTree) E(err error, message string, args ...interface{}) {
	eLogger(message, args)
}

func V(err error, message string, args ...interface{}) {
	tree_of_souls.V(err, message, args...)
}

func D(err error, message string, args ...interface{}) {
	tree_of_souls.D(err, message, args...)
}

func I(err error, message string, args ...interface{}) {
	tree_of_souls.I(err, message, args...)
}

func W(err error, message string, args ...interface{}) {
	tree_of_souls.W(err, message, args...)
}

func E(err error, message string, args ...interface{}) {
	tree_of_souls.E(err, message, args...)
}
