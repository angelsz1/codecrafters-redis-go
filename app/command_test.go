package main

import (
	"testing"
	"time"
)

const (
	NULL_BULK = "$-1\r\n"
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

func TestSetAndGetSuccesfully(t *testing.T) {
	expected := EncodeAsSimpleString("OK")
	out := ProcessComand([]string{"set", "key", "value"})
	AssertEqual(expected, out, t)
	expected = EncodeAsBulk([]string{"value"})
	out = ProcessComand([]string{"get", "key"})
	AssertEqual(expected, out, t)
}

func TestSetAndGetKeyNotFound(t *testing.T) {
	expected := EncodeAsSimpleString("OK")
	out := ProcessComand([]string{"set", "key", "value"})
	AssertEqual(expected, out, t)
	expected = NULL_BULK
	out = ProcessComand([]string{"get", "pato"})
	AssertEqual(expected, out, t)
}

func TestSetPXAndGetOK(t *testing.T) {
	expected := EncodeAsSimpleString("OK")
	out := ProcessComand([]string{"set", "key", "value", "px", "10"})
	AssertEqual(expected, out, t)
	time.Sleep(10 * time.Millisecond)
	expected = EncodeAsBulk([]string{"value"})
	out = ProcessComand([]string{"get", "key"})
	AssertEqual(expected, out, t)
}
