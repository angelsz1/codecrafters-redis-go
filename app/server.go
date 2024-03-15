package main

import (
	"fmt"
	"net"
	"os"
)

var state map[string]string = map[string]string{
	"port":        "6379",
	"role":        "master",
	"master_host": "localhost",
	"master_port": "6379",
}

func main() {
	setup()
	tcpDirection := fmt.Sprintf("0.0.0.0:%s", state["port"])
	l, err := net.Listen("tcp", tcpDirection)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	conn, err := l.Accept()
	for err == nil {
		go handleConnection(conn)
		conn, err = l.Accept()
	}
}

func setup() {
	setUpFlags()
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
