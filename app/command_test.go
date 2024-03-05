package main

import (
	"testing"
	"reflect"
)

func TestPing(t *testing.T) {
	expected := EncodeAsBulk([]string{"PONG"})
	out := EncodeAsBulk(ProcessComand([]string{"pInG"}))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}

func TestEcho(t *testing.T) {
	expected := EncodeAsBulk([]string{"helloWorld"})
	out := EncodeAsBulk(ProcessComand([]string{"echo", "helloWorld"}))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}