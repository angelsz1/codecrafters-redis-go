package main

import (
	"strconv"
	"fmt"
	"regexp"
)

const (
	Integer = ':'
	String  = '+'
	Bulk    = '$'
	Array   = '*'
	Error   = '-'
)

func ReadRESP(r []byte) []string {
	switch r[0] {
	case Integer:
		return readInteger(r)
	case String:
		return readString(r)
	case Bulk:
		return readBulk(r)
	case Array:
		return readArray(r)
	// case Error:
	// 	return readError(r)
	}
	return []string{"FATAL ERROR"}
}

func readString(r []byte) []string {
	length := len(r)
	str := r[1:length-2]
	return []string{string(str)}
}

func readInteger(r []byte) []string {
	length := len(r)
	num := r[1:length-2]
	conv, err := strconv.Atoi(string(num))
	if err != nil {
		panic("This shouldn't happen")
	}
	return []string{fmt.Sprintf("%d", conv)}
}

func readBulk(r []byte) []string {
	index := 1
	for r[index] != '\r' {
		index++
	}
	conv, err := strconv.Atoi(string(r[1:index]))
	if err != nil {
		panic("Corrupt bulk string")
	}
	if conv == 0 {
		return []string{""}
	}
	return []string{string(r[index + 2:len(r)-2])}
}

func readArray(r []byte) []string {
	var results []string
	regexStr := "\n"
	reg := regexp.MustCompile(regexStr)
	matches := reg.FindAllStringSubmatchIndex(string(r), -1)
	prev_match := -1
	for i, match := range matches {
		if prev_match == -1 {
			prev_match = match[0]
			continue
		}
		if i % 2 != 0 {
			continue
		} 
		results = append(results, ReadRESP(r[prev_match+1:match[0]+1])[0])
		prev_match = match[0]
	}
	return results
}

func EncodeAsBulk(str []string) []byte {
	if len(str) != 1 {
		return nil
	}
	actStr := str[0]
	strLen := len(actStr)
	encodedStr := fmt.Sprintf("$%d\r\n%s\r\n", strLen, actStr)
	return []byte(encodedStr)
}