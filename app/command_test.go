package main

import (
	"testing"
	"reflect"
)

func TestPing(t *testing.T) {
	expected := EncodeAsBulk([]string{"pong"})
	out := EncodeAsBulk(ProcessComand([]string{"ping"}))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}