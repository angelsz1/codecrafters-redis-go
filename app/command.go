package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Expiry struct {
	deadline int64
	setTime  int64
}

var registry map[string]func([]string) string = map[string]func([]string) string{
	"ping":     ping,
	"echo":     echo,
	"set":      set,
	"get":      get,
	"info":     info,
	"replconf": replconf,
	"psync":    psync,
}

var values map[string]string = make(map[string]string)
var expiryValues map[string]Expiry = make(map[string]Expiry)

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
	if len(cmd) > 3 {
		pxSet(cmd)
	}
	values[cmd[1]] = cmd[2]
	if state["role"] == "master" {
		return EncodeAsSimpleString("OK")
	}
	return ""
}

func pxSet(cmd []string) {
	conv, err := strconv.Atoi(cmd[4])
	if err != nil {
		panic("Error: wrong expiry")
	}
	expiryValues[cmd[1]] = Expiry{int64(conv), time.Now().UnixMilli()}
}

func get(cmd []string) string {
	_, ok := values[cmd[1]]
	if ok {
		_, ok := expiryValues[cmd[1]]
		if !ok {
			return EncodeAsBulk([]string{values[cmd[1]]})
		} else {
			if time.Now().UnixMilli()-expiryValues[cmd[1]].setTime <= expiryValues[cmd[1]].deadline {
				return EncodeAsBulk([]string{values[cmd[1]]})
			} else {
				deleteKeyValue(cmd[1])
			}
		}
	}
	return EncodeAsBulk([]string{"null"})
}

func info(cmd []string) string {
	if strings.Compare("replication", cmd[1]) == 0 {
		return EncodeAsBulk([]string{fmt.Sprintf("role:%s\nmaster_replid:%s\nmaster_repl_offset:%s", state["role"], state["replication_id"], state["replication_offset"])})
	}
	return EncodeAsBulk([]string{"null"})
}

func deleteKeyValue(key string) {
	delete(values, key)
	delete(expiryValues, key)
}

func replconf(cmd []string) string {
	return EncodeAsSimpleString("OK")
}

func psync(cmd []string) string {
	fullRsync := "FULLRESYNC"
	replId := state["replication_id"]
	replOff := state["replication_offset"]
	return EncodeAsSimpleString(fmt.Sprintf("%s %s %s", fullRsync, replId, replOff)) +
		RDBState()
}
