package main

import "reflect"

// probably not idiomatic Go, but a fun way to do verbose/debug only print statements

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
