package main

import (
	"fmt"
	"net"
	"os"
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
	if state["role"] == "master" {
		l := connectToHost("0.0.0.0", state["port"])
		connectionHandler(l)
		//state["replication_id"] = RandomString()
	} else {
		//repl handshake
		sendHandshake()
	}
}

func sendHandshake() {
	l := connectToHost(state["master_host"], state["master_port"])
	pingMaster(l)
	//two more steps
}

func pingMaster(l net.Listener) {
	conn, err := l.Accept()
	if err == nil {
		wBuf := EncodeAsBulkArray([]string{"ping"})
		conn.Write([]byte(wBuf))
	}
}

func setUpFlags() {
	args := os.Args[1:]
	for idx, value := range args {
		switch value {
		case "--port":
			state["port"] = args[idx+1]
		case "--replicaof":
			state["role"] = "slave"
			state["master_host"] = args[idx+1]
			state["master_port"] = args[idx+2]
		}
	}
}

func connectToHost(host string, port string) net.Listener {
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	defer l.Close()
	return l
}

func connectionHandler(l net.Listener) {
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
