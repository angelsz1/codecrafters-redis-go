package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type replica struct {
	conn      net.Conn
	bytesSent int
}

var WriteCommands = []string{
	"set",
	"del",
	"replconf",
}

const RDBMasterStateFilePath = "empty.rdb"

var replicas []replica

func RDBState() string {
	data, err := os.ReadFile(RDBMasterStateFilePath)
	if err != nil {
		panic("RDB file corrupted")
	}
	return formatRDB(string(data))
}

func formatRDB(data string) string {
	formattedStr := "$"
	formattedStr += fmt.Sprintf("%d\r\n", len(data))
	formattedStr += data
	return formattedStr
}

func AddReplica(conn net.Conn) {
	replicas = append(replicas, replica{conn, 0})
}

func Propagate(command []byte) {
	strCmd := strings.ReplaceAll(string(command), "\x00", "")
	for i, replica := range replicas {
		writer := bufio.NewWriter(replica.conn)
		writer.WriteString(strCmd)
		writer.Flush()
		replicas[i].bytesSent += len(strCmd)
	}
}

func IsWriteCommand(command []string) bool {
	if len(command) <= 0 {
		return false
	}
	cmdName := strings.ToLower(command[0])
	for _, cmd := range WriteCommands {
		if strings.Compare(cmd, cmdName) == 0 {
			return true
		}
	}
	return false
}

func ReplicaExists(conn net.Conn) bool {
	for _, repl := range replicas {
		if conn == repl.conn {
			return true
		}
	}
	return false
}

func propagateToReplica(cmd []byte, repl replica) {
	writer := bufio.NewWriter(repl.conn)
	writer.Write(cmd)
	writer.Flush()
}
