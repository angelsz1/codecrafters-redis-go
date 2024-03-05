package main

import (
	"testing"
)

func TestReadString(t *testing.T) {
	expected := []string{"Hola"}
	out := ReadRESP([]byte("+Hola\r\n"))
	AssertEqual(expected, out, t)
}

func TestReadInteger(t *testing.T) {
	expected := []string{"150"}
	out := ReadRESP([]byte(":+150\r\n"))
	AssertEqual(expected, out, t)
}

func TestReadBulkString(t *testing.T) {
	expected := []string{"Hola"}
	out := ReadRESP([]byte("$5\r\nHola\r\n"))
	AssertEqual(expected, out, t)
}

func TestReadBulkStringEmpty(t *testing.T) {
	expected := []string{""}
	out := ReadRESP([]byte("$0\r\n\r\n"))
	AssertEqual(expected, out, t)
}

func TestReadArray(t *testing.T) {
	expected := []string{"Hola", "Bob", "Esponja"}
	out := ReadRESP([]byte("*3\r\n$4\r\nHola\r\n$3\r\nBob\r\n$7\r\nEsponja\r\n"))
	AssertEqual(expected, out, t)
}

func TestBulkEncoding(t *testing.T) {
	expected := "$3\r\nhey\r\n"
	out := EncodeAsBulk([]string{"hey"})
	AssertEqual(expected, out, t)
}

func TestBulkEncoding2(t *testing.T) {
	expected := "$-1\r\n"
	out := EncodeAsBulk([]string{"null"})
	AssertEqual(expected, out, t)
}

func TestEncodeSimpleString(t *testing.T) {
	expected := "+OK\r\n"
	out := EncodeAsSimpleString("OK")
	AssertEqual(expected, out, t)
}
