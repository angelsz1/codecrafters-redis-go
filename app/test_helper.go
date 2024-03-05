package main

import (
	"testing"
	"reflect"
)

func AssertEqual(expected interface{}, out interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}