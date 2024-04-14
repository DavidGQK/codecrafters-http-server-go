package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

const (
	RESPONSE_200 = "HTTP/1.1 200 OK\r\n"
	RESPONSE_201 = "HTTP/1.1 201 OK\r\n"
	RESPONSE_404 = "HTTP/1.1 404 Not Found\r\n"
	SEP          = "\r\n"
	DOUBLE_SEP   = "\r\n\r\n"
)

const (
	GET_METHOD  = "get"
	POST_METHOD = "post"
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
	body            []byte
}

func parseReq(query []byte) request {
	tokens := bytes.Split(query, []byte("\r\n"))
	reqInfo := bytes.Split(tokens[0], []byte(" "))

	parsedReq := request{
		method:          strings.ToLower(string(reqInfo[0])),
		URI:             string(reqInfo[1]),
		protocolVersion: strings.ToLower(string(reqInfo[2])),
		headers:         make(map[string]string),
	}

	bodyPos := 0
	for i, token := range tokens[1:] {
		if len(token) == 0 {
			bodyPos = i + 2

			if parsedReq.method == POST_METHOD {
				parsedReq.body = bytes.Join(tokens[bodyPos:], []byte("\r\n"))
			}

			break
		}

		header := strings.SplitN(string(token), ":", 2)
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

	if req.method == GET_METHOD {
		if req.URI == PATH_EMPTY {
			resp.WriteString(RESPONSE_200 + SEP)
		} else {
			switch pathType {

			case PATH_ECHO:
				reqData := tokens[1]
				resp.WriteString(RESPONSE_200)
				resp.WriteString("Content-Type: text/plain" + SEP)
				resp.WriteString("Content-Length: ")
				resp.WriteString(fmt.Sprint(len(reqData)))
				resp.WriteString(DOUBLE_SEP)
				resp.WriteString(reqData)
				resp.WriteString(DOUBLE_SEP)

			case PATH_USER_AGENT:
				data := req.headers[PATH_USER_AGENT]
				resp.WriteString(RESPONSE_200)
				resp.WriteString("Content-Type: text/plain" + SEP)
				resp.WriteString("Content-Length: ")
				resp.WriteString(fmt.Sprint(len(data)))
				resp.WriteString(DOUBLE_SEP)
				resp.WriteString(data)
				resp.WriteString(DOUBLE_SEP)

			case PATH_FILES:
				filename := tokens[1]
				fileContent, fileExist := findFile(conf.directory, filename)

				if fileExist && conf.dirExists {
					resp.WriteString(RESPONSE_200)
					resp.WriteString("Content-Type: application/octet-stream" + SEP)
					resp.WriteString("Content-Length: ")
					resp.WriteString(fmt.Sprint(len(fileContent)))
					resp.WriteString(DOUBLE_SEP)
					resp.WriteString(fileContent)
					resp.WriteString(DOUBLE_SEP)
				} else {
					resp.WriteString(RESPONSE_404 + SEP)
				}

			default:
				resp.WriteString(RESPONSE_404 + SEP)
			}
		}
	}

	if req.method == POST_METHOD {
		switch pathType {
		case PATH_FILES:
			filename := tokens[1]

			err := saveFile(conf.directory, filename, req.body)
			if err != nil {
				fmt.Println("Error writing while saving file")
				resp.WriteString(RESPONSE_404 + SEP)
				break
			}

			resp.WriteString(RESPONSE_201 + SEP)

		default:
			resp.WriteString(RESPONSE_404 + SEP)
		}
	}

	_, err := conn.Write([]byte(resp.String()))
	if err != nil {
		fmt.Println("Error answering request:", err.Error())
	}

	fmt.Println(resp.String())
}
