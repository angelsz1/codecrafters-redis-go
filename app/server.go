package main

import (
	"fmt"
	"net"
	"os"
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
		SendHandshake()
	}
	l := connectToHost(HOST, state["port"])
	connectionHandler(l)
}

func setUpFlags() {
	args := os.Args[1:]
	for idx, value := range args {
		switch value {
		case "--port":
			state["port"] = args[idx+1]
		case "--replicaof":
			SetReplicaState(args[idx+1], &state)
		}
	}
}

func connectToHost(host string, port string) net.Listener {
	l, err := net.Listen("tcp", net.JoinHostPort(host, port))
	if err != nil {
		fmt.Println("Failed to bind to port " + port)
		os.Exit(1)
	}
	return l
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
		if state["role"] == "master" && !IsHandshakeCommand(ReadRESP(rBuf)) {
			Propagate(rBuf)
		} else if state["role"] == "master" && !ReplicaExists(conn) {
			AddReplica(conn)
		}
		wBuf := ProcessComand(ReadRESP(rBuf))
		conn.Write([]byte(wBuf))
		_, err = conn.Read(rBuf)
	}
}
