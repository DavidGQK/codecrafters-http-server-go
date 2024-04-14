package main

import (
	"fmt"
	"net"
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
	PATH_EMPTY      = "/"
	PATH_ECHO       = "echo"
	PATH_USER_AGENT = "user-agent"
	PATH_FILES      = "files"
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
		key = strings.TrimSpace(key)
		val := header[1]
		val = strings.TrimSpace(val)

		parsedReq.headers[key] = val
	}

	return parsedReq
}

func ansReq(conn net.Conn, req request) {
	var resp strings.Builder

	path := strings.Trim(req.URI, "/")
	tokens := strings.SplitN(path, "/", 2)
	pathType := strings.ToLower(tokens[0])

	//fmt.Println(req.URI)
	//fmt.Println(fmt.Sprintf("path: %s, tokens: %s, its len: %d", path, tokens, len(tokens)))

	if req.method == GET_METHOD {
		if req.URI == PATH_EMPTY {
			resp.WriteString(RESPONSE_200)
		} else {
			switch pathType {

			case PATH_ECHO:
				fmt.Println("echo case")
				reqData := tokens[1]
				resp.WriteString("HTTP/1.1 200 OK\r\n")
				resp.WriteString("Content-Type: text/plain\r\n")
				resp.WriteString("Content-Length: ")
				resp.WriteString(fmt.Sprint(len(reqData)))
				resp.WriteString("\r\n\r\n")
				resp.WriteString(reqData)
				resp.WriteString("\r\n\r\n")

			case PATH_USER_AGENT:
				fmt.Println("user-agent case")
				data := req.headers[PATH_USER_AGENT]
				resp.WriteString("HTTP/1.1 200 OK\r\n")
				resp.WriteString("Content-Type: text/plain\r\n")
				resp.WriteString("Content-Length: ")
				resp.WriteString(fmt.Sprint(len(data)))
				resp.WriteString("\r\n\r\n")
				resp.WriteString(data)
				resp.WriteString("\r\n\r\n")

			case PATH_FILES:
				fmt.Println("files case")
				filename := tokens[1]
				fileContent, fileExist := findFile(conf.directory, filename)
				fmt.Println("filename", filename, "conf.directory", conf.directory)
				if fileExist && conf.dirExists {
					resp.WriteString("HTTP/1.1 200 OK\r\n")
					resp.WriteString("Content-Type: application/octet-stream\r\n")
					resp.WriteString("Content-Length: ")
					resp.WriteString(fmt.Sprint(len(fileContent)))
					resp.WriteString("\r\n\r\n")
					resp.WriteString(fileContent)
					resp.WriteString("\r\n\r\n")
				} else {
					resp.WriteString(RESPONSE_404)
				}

			default:
				fmt.Println("default case")
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
