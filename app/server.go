package main

import (
	"fmt"
	"net"
	"os"
)

var conf *config

func handleConnection(conn net.Conn) {

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error handling connection:", err.Error())
		return
	}

	res := parseReq(buf[:n])
	ansReq(conn, res)
	defer conn.Close()
}

func main() {
	conf = getConfig()

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		fmt.Println("Got a connection!")
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}
