package main

import (
	"strings"
)

var registry map[string]func([]string) string = map[string]func([]string) string {
	"ping": ping,
	"echo": echo,
	"set" : set,
	"get" : get,
}

var values map[string]string = make(map[string]string, 1024)

func ProcessComand(cmd []string) string {
	command := strings.ToLower(cmd[0])
	_, ok := registry[command]
	if ok {
		return registry[command](cmd)
	}
	return ""
}

func ping(cmd []string) string {
	return EncodeAsBulk([]string{"PONG"})
}

func echo(cmd []string) string {
	return EncodeAsBulk([]string{cmd[1]})
}

func set(cmd []string) string {
	values[cmd[1]] = cmd[2]
	return EncodeAsSimpleString("OK")
}

func get(cmd []string) string {
	_, ok := values[cmd[1]]
	if ok {
		return EncodeAsBulk([]string{values[cmd[1]]})
	}
	return EncodeAsBulk([]string{"null"})
}
