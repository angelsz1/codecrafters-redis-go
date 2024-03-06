package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	portFlag := flag.Int("port", 6379, "redis connection port")
	flag.Parse()
	tcpDirection := fmt.Sprintf("0.0.0.0:%d", *portFlag)
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
