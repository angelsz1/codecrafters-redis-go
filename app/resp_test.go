package main

import (
	"testing"
	"reflect"
)

func TestReadString(t *testing.T) {
	expected := []string{"Hola"}
	out := ReadRESP([]byte("+Hola\r\n"))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}

func TestReadInteger(t *testing.T) {
	expected := []string{"150"}
	out := ReadRESP([]byte(":+150\r\n"))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}

func TestReadBulkString(t *testing.T) {
	expected := []string{"Hola"}
	out := ReadRESP([]byte("$5\r\nHola\r\n"))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}

func TestReadBulkStringEmpty(t *testing.T) {
	expected := []string{""}
	out := ReadRESP([]byte("$0\r\n\r\n"))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}

func TestReadArray(t *testing.T) {
	expected := []string{"Hola", "Bob", "Esponja"}
	out := ReadRESP([]byte("*3\r\n$4\r\nHola\r\n$3\r\nBob\r\n$7\r\nEsponja\r\n"))
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}

func TestBulkEncoding(t *testing.T) {
	expected := []byte("$3\r\nhey\r\n")
	out := EncodeAsBulk([]string{"hey"})
	if !reflect.DeepEqual(expected, out) {
		t.Errorf("Not equal: \nExpected -> %s\nOut -> %s", expected, out)
	}
}