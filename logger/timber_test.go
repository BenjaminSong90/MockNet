package logger

import (
	"testing"
)

type TestTree struct {
	writeInfo string
}

func (t *TestTree) V(message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) D(message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) I(message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) W(message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) E(message string, args ...interface{}) {
	t.writeInfo = message
}

func TestTimber(t *testing.T) {
	testTree := &TestTree{}
	PlantTree(testTree)

	I("test")

	if testTree.writeInfo != "test" {
		t.Error()
	}

	if Size() != 1 {
		t.Error()
	}

	PlantTree(&DebugTree{})

	if Size() != 2 {
		t.Error()
	}

	UnRoot(testTree)

	if Size() != 1 {
		t.Error()
	}

	UnRootAll()

	if Size() != 0 {
		t.Error()
	}
}
