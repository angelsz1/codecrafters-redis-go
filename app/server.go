package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println(ReadRESP([]byte("*3\r\n$4\r\nHola\r\n$3\r\nBob\r\n$7\r\nEsponja\r\n")))
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	conn, err := l.Accept()
	for err == nil {
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
		conn, err = l.Accept()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	rBuf := make([]byte, 1024)
	_, err := conn.Read(rBuf)
	for err == nil {
		wBuf := ReadRESP(rBuf)
		conn.Write(EncodeAsBulk(wBuf))
		_, err = conn.Read(rBuf)
	}
}
