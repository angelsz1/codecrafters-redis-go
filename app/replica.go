package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func SendHandshake() {
	l := connectToMaster()
	pingMaster(l)
	buffer := make([]byte, 1024)
	l.Read(buffer)
	replconfMaster(l)
	l.Read(buffer)
	psyncMaster(l)
	//one more steps
}

func pingMaster(l net.Conn) {
	wBuf := EncodeAsBulkArray([]string{"ping"})
	l.Write([]byte(wBuf))
}

func replconfMaster(l net.Conn) {
	wBuf := EncodeAsBulkArray([]string{"REPLCONF", "listening-port", state["port"]})
	l.Write([]byte(wBuf))
	wBuf = EncodeAsBulkArray([]string{"REPLCONF", "capa", "psync2"})
	l.Write([]byte(wBuf))
}

func psyncMaster(l net.Conn) {
	//for now, it looks hardcoded
	buffer := make([]byte, 1024)
	l.Read(buffer)
	wBuf := EncodeAsBulkArray([]string{"PSYNC", "?", "-1"})
	l.Write([]byte(wBuf))
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
