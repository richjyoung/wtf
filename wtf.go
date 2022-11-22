// WTF Is This?
//
// Golang is great until at runtime you are dealing with an interface for which the underlying implementation is buried in a package somewhere.
// Yes it's easier than other languages to click through and find it, or...
//
//	fmt.Println(wtf.IsThis(thing))
//
// See the code in wtf_test.go for examples, here are some excerpts
//
//	aStruct := wtf.IsThis(TestStruct{})
//	assert.Equal(t, "github.com/richjyoung/wtf_test.TestStruct", aStruct)
//
//	aStructPointer := wtf.IsThis(ptr)
//	assert.Equal(t, "*github.com/richjyoung/wtf_test.TestStruct", aStructPointer)
//
//	aStructPointerPointer := wtf.IsThis(&ptr)
//	assert.Equal(t, "**github.com/richjyoung/wtf_test.TestStruct", aStructPointerPointer)
//
//	anArray := wtf.IsThis([1]TestStruct{})
//	assert.Equal(t, "[1]github.com/richjyoung/wtf_test.TestStruct", anArray)
//
//	anArrayOfPointers := wtf.IsThis([1]*TestStruct{})
//	assert.Equal(t, "[1]*github.com/richjyoung/wtf_test.TestStruct", anArrayOfPointers)
//
//	aChan := wtf.IsThis(ch)
//	assert.Equal(t, "chan github.com/richjyoung/wtf_test.TestStruct", aChan)
//
//	anRChan := wtf.IsThis(rch)
//	assert.Equal(t, "<-chan github.com/richjyoung/wtf_test.TestStruct", anRChan)
//
//	anSChan := wtf.IsThis(sch)
//	assert.Equal(t, "chan<- github.com/richjyoung/wtf_test.TestStruct", anSChan)
//
//	aSlice := wtf.IsThis([]TestStruct{})
//	assert.Equal(t, "[]github.com/richjyoung/wtf_test.TestStruct", aSlice)
//
//	aSliceOfPointers := wtf.IsThis([]*TestStruct{})
//	assert.Equal(t, "[]*github.com/richjyoung/wtf_test.TestStruct", aSliceOfPointers)
//
//	aMap := wtf.IsThis(map[string]TestStruct{})
//	assert.Equal(t, "map[string]github.com/richjyoung/wtf_test.TestStruct", aMap)
//
//	aMapOfPointers := wtf.IsThis(map[string]*TestStruct{})
//	assert.Equal(t, "map[string]*github.com/richjyoung/wtf_test.TestStruct", aMapOfPointers)
//
//	aFunc := wtf.IsThis(fn)
//	assert.Equal(t, "func (*github.com/richjyoung/wtf_test.TestStruct) error {}", aFunc)
//
//	aFuncInterfaceArg := wtf.IsThis(fni)
//	assert.Equal(t, "func (github.com/richjyoung/wtf_test.TestIface) *github.com/richjyoung/wtf_test.TestStruct {}", aFuncInterfaceArg)
package wtf

import (
	"fmt"
	"reflect"
)

const NoIdea = `¯\_(ツ)_/¯`

var chanMap = map[reflect.ChanDir]string{
	reflect.RecvDir: "<-chan",
	reflect.BothDir: "chan",
	reflect.SendDir: "chan<-",
}

// wtf.IsThis returns a string representation of the target.
// It is not intended to exactly match Golang syntax, but should help work out what it needs to be.
//
// Needs a concrete value to return a valid string, empty interfaces or nil pointers will return wtf.NoIdea.
func IsThis(target interface{}) string {
	typ := reflect.TypeOf(target)
	if typ != nil {
		res := wtfIsThis(typ)
		if res != "" {
			return res
		}
	}
	return NoIdea
}

func wtfIsThis(t reflect.Type) string {
	switch t.Kind() {
	case reflect.Pointer:
		return aPointerOf(t)
	case reflect.Map:
		return aMapOf(t)
	case reflect.Slice:
		return aSliceOf(t)
	case reflect.Array:
		return anArrayOf(t)
	case reflect.Chan:
		return aChanOf(t)
	case reflect.Func:
		return aFuncOf(t)
	default:
		if t.PkgPath() != "" {
			return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
		} else if t.Name() != "" {
			return t.Name()
		}
		return NoIdea
	}
}

func aPointerOf(t reflect.Type) string { return fmt.Sprintf("*%s", wtfIsThis(t.Elem())) }
func anArrayOf(t reflect.Type) string  { return fmt.Sprintf("[%d]%s", t.Len(), wtfIsThis(t.Elem())) }
func aSliceOf(t reflect.Type) string   { return fmt.Sprintf("[]%s", wtfIsThis(t.Elem())) }

func aMapOf(t reflect.Type) string {
	return fmt.Sprintf("map[%s]%s", wtfIsThis(t.Key()), wtfIsThis(t.Elem()))
}

func aChanOf(t reflect.Type) string {
	return fmt.Sprintf("%s %s", chanMap[t.ChanDir()], wtfIsThis(t.Elem()))
}

func aFuncOf(t reflect.Type) string {
	res := "func"
	funcName := t.Name()
	if funcName != "" {
		res = fmt.Sprintf("%s %s(", res, funcName)
	} else {
		res = fmt.Sprintf("%s (", res)
	}

	for i := 0; i < t.NumIn(); i++ {
		if i > 0 {
			res = fmt.Sprintf("%s, %s", res, wtfIsThis(t.In(i)))
		} else {
			res = fmt.Sprintf("%s%s", res, wtfIsThis(t.In(i)))
		}
	}

	res = fmt.Sprintf("%s)", res)

	if t.NumOut() > 0 {
		if t.NumOut() > 1 {
			for i := 0; i < t.NumOut(); i++ {
				if i > 0 {
					res = fmt.Sprintf("%s, %s", res, wtfIsThis(t.Out(i)))
				} else {
					res = fmt.Sprintf("%s (%s", res, wtfIsThis(t.Out(i)))
				}
			}
			res = fmt.Sprintf("%s)", res)
		} else {
			res = fmt.Sprintf("%s %s", res, wtfIsThis(t.Out(0)))
		}
	}

	res = fmt.Sprintf("%s {}", res)

	return res
}