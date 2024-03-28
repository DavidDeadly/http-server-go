package handlers

import (
	"fmt"
	"io"
	"net"
	"regexp"

	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

const (
	OK          = 200
	CREATED     = 201
	NOT_FOUND   = 404
	BAD_REQUEST = 400
)

const (
	GET  utils.METHOD = "GET"
	POST utils.METHOD = "POST"
)

func HandleConnection(conn net.Conn) {
	fmt.Println("Handling...")

	defer conn.Close()

	rawRequest := make([]byte, 1024)
	bytes, err := conn.Read(rawRequest)

	if err == io.EOF {
		return
	}

	if err != nil {
		fmt.Println("Error reading data from the conn: ", err.Error())
		return
	}

	request := utils.ParseHTTPRequest(rawRequest, bytes)

	var response []byte

	switch {
	case request.Path == "/":
		response = Response(nil, OK, nil)

	case regexp.MustCompile("/echo/*").MatchString(request.Path):
		response = Echo(request)

	case request.Path == "/user-agent":
		response = UserAgent(request)

	case regexp.MustCompile("/files/*").MatchString(request.Path):
    response = Files(request)

	default:
		response = NotFoundResponse()
	}

	bytes, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error sending data to the connection: ", err.Error())
		return
	}

	fmt.Printf("Send %v bytes\n", bytes)
}

func Response(body *string, code uint16, contentType *string) []byte {
	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n", code)

	if contentType == nil {
		content := "text/plain"
		contentType = &content
	}

	if body != nil {
		response += fmt.Sprintf("Content-Type: %s\r\n", *contentType)
		response += fmt.Sprintf("Content-Length: %d\r\n", len(*body))
		response += fmt.Sprintf("\r\n%s", *body)
	}

	response += "\r\n"

	fmt.Println(response)
	return []byte(response)
}

func NotFoundResponse() []byte {
	response := fmt.Sprintf("HTTP/1.1 %d Not Found\r\n\r\n", NOT_FOUND)

	return []byte(response)
}
