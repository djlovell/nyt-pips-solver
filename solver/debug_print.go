package solver

import "reflect"

// probably not idiomatic Go, but a fun way to do verbose/debug only print statements

// TODO: convert to environment variable?
var debug bool = true

func debugPrint(printStatement any, args ...any) {
	if debug {
		f := reflect.ValueOf(printStatement)
		aList := []reflect.Value{}
		for _, a := range args {
			aList = append(aList, reflect.ValueOf(a))
		}
		f.Call(aList)
	}
}
