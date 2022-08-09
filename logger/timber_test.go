package logger

import (
	"testing"
)

type TestTree struct {
	writeInfo string
}

func (t *TestTree) V(err error, message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) D(err error, message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) I(err error, message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) W(err error, message string, args ...interface{}) {
	t.writeInfo = message
}

func (t *TestTree) E(err error, message string, args ...interface{}) {
	t.writeInfo = message
}

func TestTimber(t *testing.T) {
	testTree := &TestTree{}
	PlantTree(testTree)

	I(nil, "test")

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
