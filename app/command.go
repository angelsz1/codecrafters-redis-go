package main

import (
	"strings"
)

var responses map[string]string = map[string]string {
	"ping" : "PONG",
}
func ProcessComand(cmd []string) []string {
	command := strings.ToLower(cmd[0])
	_, ok := responses[command]
	if ok {
		return []string{responses[command]}
	}
	return []string{""}
}