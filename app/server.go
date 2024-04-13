package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

const (
	RESPONSE_200 = "HTTP/1.1 200 OK\r\n\r\n"
	RESPONSE_404 = "HTTP/1.1 404 Not Found\r\n\r\n"
)

type request struct {
	method          string
	URI             string
	protocolVersion string
	headers         map[string]string
}

func parseReq(query []byte) request {
	tokens := strings.Split(string(query), "\r\n")

	reqInfo := strings.Split(tokens[0], " ")

	parsedReq := request{
		method:          reqInfo[0],
		URI:             reqInfo[1],
		protocolVersion: reqInfo[2],
		headers:         make(map[string]string),
	}

	for _, token := range tokens[1:] {
		if token == "" {
			break
		}
		header := strings.SplitN(token, ":", 2)
		fmt.Println(header)
		parsedReq.headers[header[0]] = header[1]
	}

	return parsedReq
}

func ansReq(conn net.Conn, req request) {
	var resp string
	if req.URI == "/" {
		resp = RESPONSE_200
	} else {
		resp = RESPONSE_404
	}

	_, err := conn.Write([]byte(resp))
	if err != nil {
		fmt.Println("Error answering request:", err.Error())
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Println("Error handling connection:", err.Error())
			return
		}

		res := parseReq(buf[:n])
		ansReq(conn, res)
	}
}

func main() {
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
