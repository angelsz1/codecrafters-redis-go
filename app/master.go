package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const RDBMasterStateFilePath = "empty.rdb"

var replicas []net.Conn

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
	replicas = append(replicas, conn)
}

func Propagate(command []byte) {
	for _, conn := range replicas {
		strCmd := strings.Replace(string(command), "\x00", "", -1)
		conn.Write([]byte(strCmd))
	}
}

func ReplicaExists(conn net.Conn) bool {
	for _, repl := range replicas {
		if conn == repl {
			return true
		}
	}
	return false
}
