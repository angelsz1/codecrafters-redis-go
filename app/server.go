package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	HOST = "0.0.0.0"
)

var state map[string]string = map[string]string{
	"port":               "6379",
	"role":               "master",
	"master_host":        "localhost",
	"master_port":        "6379",
	"replication_id":     "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
	"replication_offset": "0",
}

func main() {
	setup()
}

func setup() {
	setUpFlags()
	if state["role"] == "slave" {
		//handshake
		sendHandshake()
	}
	l := connectToHost(HOST, state["port"])
	connectionHandler(l)
}

func sendHandshake() {
	l := connectToMaster()
	pingMaster(l)
	//two more steps
}

func pingMaster(l net.Conn) {
	defer l.Close()
	wBuf := EncodeAsBulkArray([]string{"ping"})
	l.Write([]byte(wBuf))
}

func setUpFlags() {
	args := os.Args[1:]
	for idx, value := range args {
		switch value {
		case "--port":
			state["port"] = args[idx+1]
		case "--replicaof":
			setReplicaState(args[idx+1])
		}
	}
}

func setReplicaState(replState string) {
	state["role"] = "slave"
	splState := strings.Split(replState, " ")
	state["master_host"] = splState[0]
	state["master_port"] = splState[1]
}

func connectToHost(host string, port string) net.Listener {
	l, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	return l
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

func connectionHandler(l net.Listener) {
	defer l.Close()
	conn, err := l.Accept()
	for err == nil {
		go handleConnection(conn)
		conn, err = l.Accept()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	rBuf := make([]byte, 1024)
	_, err := conn.Read(rBuf)
	for err == nil {
		wBuf := ProcessComand(ReadRESP(rBuf))
		conn.Write([]byte(wBuf))
		_, err = conn.Read(rBuf)
	}
}
