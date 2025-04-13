package main

import (
	"bufio"
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
	"dir":                "/tmp/redis-data",
	"dbfilename":         "dump.rdb",
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

		case "--dir":
			state["dir"] = args[idx+1]
		case "--dbfilename":
			state["dbfilename"] = args[idx+1]
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
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	for {
		line := make([]byte, 1024)
		_, err := reader.Read(line)
		if err != nil {
			break
		}
		resp := ReadRESP(line)
		if state["role"] == "master" && !IsHandshakeCommand(resp) && IsWriteCommand(resp) {
			Propagate(line)
		} else if state["role"] == "master" && !ReplicaExists(conn) && IsPsyncCommand(resp) {
			AddReplica(conn)
		}
		wBuf := ProcessComand(resp)
		writer.Write([]byte(wBuf))
		writer.Flush()
	}
}
