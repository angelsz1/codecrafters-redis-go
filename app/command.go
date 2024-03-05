package main

var responses map[string]string = map[string]string {
	"ping" : "pong",
}
func ProcessComand(cmd []string) []string {
	_, ok := responses[cmd[0]]
	if ok {
		return []string{responses[cmd[0]]}
	}
	return []string{""}
}