package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var byteCount int = 0
var canCountBytes bool = false

func SendHandshake() {
	conn := connectToMaster()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	pingMaster(writer)
	readBuf(reader)
	replconfMaster(writer)
	readBuf(reader)
	psyncMaster(reader, writer)
	go waitForMaster(conn)
}

func waitForMaster(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		line := make([]byte, 1024)
		_, err := reader.Read(line)
		if err != nil {
			break
		}
		commands := CheckForMultipleCommand(line)
		for _, cmd := range commands {
			if len(cmd) == 0 {
				continue
			}
			if IsWriteCommand(cmd) {
				wBuf := RespStringToRespArray(ProcessComand(cmd))
				writer.Write([]byte(wBuf))
				writer.Flush()
			}
			countBytes(cmd)
		}
	}
	conn.Close()
}

func countBytes(cmd []string) {
	if canCountBytes {
		byteCount += len(EncodeAsBulkArray(cmd))
	}
}

func pingMaster(writer *bufio.Writer) {
	wBuf := EncodeAsBulkArray([]string{"ping"})
	writer.Write([]byte(wBuf))
	writer.Flush()
}

func replconfMaster(writer *bufio.Writer) {
	writer.Write([]byte(EncodeAsBulkArray([]string{"REPLCONF", "listening-port", state["port"]})))
	writer.Write([]byte(EncodeAsBulkArray([]string{"REPLCONF", "capa", "psync2"})))
	writer.Flush()
}

func psyncMaster(reader *bufio.Reader, writer *bufio.Writer) {
	readBuf(reader)
	writer.Write([]byte(EncodeAsBulkArray([]string{"PSYNC", "?", "-1"})))
	writer.Flush()
}

func readBuf(reader *bufio.Reader) {
	// Reads until newline for simplicity, adjust if protocol differs
	reader.ReadBytes('\n')
}

func SetReplicaState(replState string, st *map[string]string) {
	(*st)["role"] = "slave"
	splState := strings.Split(replState, " ")
	(*st)["master_host"] = splState[0]
	(*st)["master_port"] = splState[1]
}

func connectToMaster() net.Conn {
	address := net.JoinHostPort(state["master_host"], state["master_port"])
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Failed to connect to master")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return conn
}

func IsHandshakeCommand(command []string) bool {
	lcmd := strings.ToLower(command[0])
	return lcmd == "ping" || lcmd == "replconf" || lcmd == "psync"
}

func IsPsyncCommand(command []string) bool {
	lcmd := strings.ToLower(command[0])
	return lcmd == "psync"
}
