package main

import (
	"testing"
)

func TestPing(t *testing.T) {
	expected := EncodeAsBulk([]string{"PONG"})
	out := ProcessComand([]string{"pInG"})
	AssertEqual(expected, out, t)
}

func TestEcho(t *testing.T) {
	expected := EncodeAsBulk([]string{"helloWorld"})
	out := ProcessComand([]string{"echo", "helloWorld"})
	AssertEqual(expected, out, t)
}

func TestSet(t *testing.T) {
	expected := EncodeAsSimpleString("OK")
	out := ProcessComand([]string{"set", "key", "value"})
	AssertEqual(expected, out, t)
}

func TestGet(t *testing.T) {
	expected := EncodeAsBulk([]string{"value"})
	out := ProcessComand([]string{"get", "value"})
	AssertEqual(expected, out, t)
}
