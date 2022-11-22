// wtf.IsThis(thing) either gives as full a string representation
// of a type as possible, or wtf.NoIdea
package wtf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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

// wtf.IsThisError returns a string representation of the chain of errors.
// Each line contains wtf.IsThis for the error type, as well as the message output of the Error() interface.
//
// If invoked with nil, the result is wtf.NoIdea
func IsThisError(err error) string {
	e := err
	str := ""
	i := 0
	for e != nil {
		if i > 0 {
			str = fmt.Sprintf("%s\n", str)
		}
		str = fmt.Sprintf("%s%s%s[%s]", str, strings.Repeat("  ", i), IsThis(e), e.Error())
		e = errors.Unwrap(e)
		i++
	}
	if str != "" {
		return str
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
