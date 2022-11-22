package wtf_test

import (
	"fmt"
	"testing"

	"github.com/richjyoung/wtf"
	"github.com/stretchr/testify/assert"
)

func TestAllTheThings(t *testing.T) {
	type TestStruct struct{}
	type TestIface interface{}

	ptr := &TestStruct{}
	var i interface{}
	var ti TestIface
	ch := make(chan TestStruct, 1)
	rch := make(<-chan TestStruct, 1)
	sch := make(chan<- TestStruct, 1)
	fn := func(*TestStruct) error { return nil }
	fni := func(TestIface) *TestStruct { return nil }

	assert.Equal(t, wtf.NoIdea, wtf.IsThis(nil))

	anInt := wtf.IsThis(0)
	assert.Equal(t, "int", anInt)

	anInterface := wtf.IsThis(i)
	assert.Equal(t, wtf.NoIdea, anInterface)

	anotherInterface := wtf.IsThis(ti)
	assert.Equal(t, wtf.NoIdea, anotherInterface)

	aStruct := wtf.IsThis(TestStruct{})
	assert.Equal(t, "github.com/richjyoung/wtf_test.TestStruct", aStruct)

	aStructPointer := wtf.IsThis(ptr)
	assert.Equal(t, "*github.com/richjyoung/wtf_test.TestStruct", aStructPointer)

	aStructPointerPointer := wtf.IsThis(&ptr)
	assert.Equal(t, "**github.com/richjyoung/wtf_test.TestStruct", aStructPointerPointer)

	anArray := wtf.IsThis([1]TestStruct{})
	assert.Equal(t, "[1]github.com/richjyoung/wtf_test.TestStruct", anArray)

	anArrayOfPointers := wtf.IsThis([1]*TestStruct{})
	assert.Equal(t, "[1]*github.com/richjyoung/wtf_test.TestStruct", anArrayOfPointers)

	aChan := wtf.IsThis(ch)
	assert.Equal(t, "chan github.com/richjyoung/wtf_test.TestStruct", aChan)

	anRChan := wtf.IsThis(rch)
	assert.Equal(t, "<-chan github.com/richjyoung/wtf_test.TestStruct", anRChan)

	anSChan := wtf.IsThis(sch)
	assert.Equal(t, "chan<- github.com/richjyoung/wtf_test.TestStruct", anSChan)

	aSlice := wtf.IsThis([]TestStruct{})
	assert.Equal(t, "[]github.com/richjyoung/wtf_test.TestStruct", aSlice)

	aSliceOfPointers := wtf.IsThis([]*TestStruct{})
	assert.Equal(t, "[]*github.com/richjyoung/wtf_test.TestStruct", aSliceOfPointers)

	aMap := wtf.IsThis(map[string]TestStruct{})
	assert.Equal(t, "map[string]github.com/richjyoung/wtf_test.TestStruct", aMap)

	aMapOfPointers := wtf.IsThis(map[string]*TestStruct{})
	assert.Equal(t, "map[string]*github.com/richjyoung/wtf_test.TestStruct", aMapOfPointers)

	aFunc := wtf.IsThis(fn)
	assert.Equal(t, "func (*github.com/richjyoung/wtf_test.TestStruct) error {}", aFunc)

	aFuncInterfaceArg := wtf.IsThis(fni)
	assert.Equal(t, "func (github.com/richjyoung/wtf_test.TestIface) *github.com/richjyoung/wtf_test.TestStruct {}", aFuncInterfaceArg)
}

func TestError(t *testing.T) {
	e1 := fmt.Errorf("error 1")
	e2 := fmt.Errorf("error 2 - %w", e1)
	e3 := fmt.Errorf("error 3 - %w", e2)
	e4 := fmt.Errorf("error 4 - %w", e3)

	assert.Equal(
		t,
		`*fmt.wrapError[error 4 - error 3 - error 2 - error 1]
  *fmt.wrapError[error 3 - error 2 - error 1]
    *fmt.wrapError[error 2 - error 1]
      *errors.errorString[error 1]`,
		wtf.IsThisError(e4))

	assert.Equal(t, wtf.NoIdea, wtf.IsThisError(nil))
}
