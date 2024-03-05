package main

import (
	"strings"
)

var registry map[string]func([]string) []string = map[string]func([]string) []string {
	"ping": ping,
	"echo": echo,
}

func ProcessComand(cmd []string) []string {
	command := strings.ToLower(cmd[0])
	_, ok := registry[command]
	if ok {
		return registry[command](cmd)
	}
	return []string{""}
}

func ping(cmd []string) []string {
	return []string{"PONG"}
}

func echo(cmd []string) []string {
	return []string{cmd[1]}
}
