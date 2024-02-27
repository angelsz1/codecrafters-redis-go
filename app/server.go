package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
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
		handleConnection(conn)
		conn, err = l.Accept()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	rBuf := make([]byte, 1024)
	wBuf := []byte("+PONG\r\n")
	_, err := conn.Read(rBuf)
	for err == nil {
		conn.Write(wBuf)
		_, err = conn.Read(rBuf)
	}

}
