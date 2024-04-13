package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	RESPONSE_200 = "HTTP/1.1 200 OK\r\n\r\n"
	RESPONSE_404 = "HTTP/1.1 404 Not Found\r\n\r\n"
)

const (
	GET_METHOD = "get"
)

const (
	PATH_EMPTY = "/"
	PATH_ECHO  = "echo"
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
		method:          strings.ToLower(reqInfo[0]),
		URI:             reqInfo[1],
		protocolVersion: strings.ToLower(reqInfo[2]),
		headers:         make(map[string]string),
	}

	for _, token := range tokens[1:] {
		if token == "" {
			break
		}
		header := strings.SplitN(token, ":", 2)
		key := strings.ToLower(header[0])
		val := strings.ToLower(header[1])
		parsedReq.headers[key] = val
	}

	return parsedReq
}

func ansReq(conn net.Conn, req request) {
	var resp strings.Builder

	path := strings.Trim(req.URI, "/")
	tokens := strings.SplitN(path, "/", 2)
	reqType := strings.ToLower(tokens[0])

	fmt.Println(req.URI)
	fmt.Println(fmt.Sprintf("path: %s, tokens: %s, its len: %d", path, tokens, len(tokens)))

	if req.method == GET_METHOD {
		if req.URI == PATH_EMPTY {
			resp.WriteString(RESPONSE_200)
		} else {
			switch reqType {
			case PATH_ECHO:
				reqData := tokens[1]
				resp.WriteString("HTTP/1.1 200 OK\r\n")
				resp.WriteString("Content-Type: text/plain\r\n")
				resp.WriteString("Content-Length: ")
				resp.WriteString(fmt.Sprint(len(reqData)))
				resp.WriteString("\r\n\r\n")
				resp.WriteString(reqData)
				resp.WriteString("\r\n\r\n")
			default:
				resp.WriteString(RESPONSE_404)
			}
		}
	}

	_, err := conn.Write([]byte(resp.String()))
	if err != nil {
		fmt.Println("Error answering request:", err.Error())
	}

	fmt.Println(resp.String())
}

func handleConnection(conn net.Conn) {

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		//if err == io.EOF {
		//	continue
		//}
		fmt.Println("Error handling connection:", err.Error())
		return
	}

	res := parseReq(buf[:n])
	ansReq(conn, res)
	defer conn.Close()
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
