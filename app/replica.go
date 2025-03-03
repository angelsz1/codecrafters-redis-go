package main

import (
	"net"
	"strings"
)

func SendHandshake() {
	l := connectToMaster()
	PingMaster(l)
	//two more steps
}

func PingMaster(l net.Conn) {
	defer l.Close()
	wBuf := EncodeAsBulkArray([]string{"ping"})
	l.Write([]byte(wBuf))
}

func SetReplicaState(replState string, st *map[string]string) {
	(*st)["role"] = "slave"
	splState := strings.Split(replState, " ")
	(*st)["master_host"] = splState[0]
	(*st)["master_port"] = splState[1]
}
