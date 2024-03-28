package handlers

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
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
		string := strings.Replace(request.Path, "/echo/", "", 1)

		response = Response(&string, OK, nil)

	case request.Path == "/user-agent":
		userAgent, ok := request.Headers["User-Agent"]

		if !ok {
			response = Response(nil, BAD_REQUEST, nil)
		} else {
			response = Response(&userAgent, OK, nil)
		}

	case regexp.MustCompile("/files/*").MatchString(request.Path):
		fileName := strings.Replace(request.Path, "/files/", "", 1)
		filePath := fmt.Sprintf("%s/%s", config.CONFIG[config.DIR], fileName)

		switch {
		case request.Method == GET:
			fmt.Println("Reading... ", filePath)

			data, err := os.ReadFile(filePath)

			if err == nil {
				body := string(data[0:])
				contentType := "application/octet-stream"
				response = Response(&body, OK, &contentType)
			} else {
				response = NotFoundResponse()
			}
		case request.Method == POST:
			err := os.WriteFile(filePath, []byte(request.Body), 0644)
      // TODO: add support for create files inside directories
      fmt.Println("Writing ... ", filePath, " with: ", request.Body)

			if err == nil {
        fmt.Printf("Error writing the file '%s' with: %s", fileName, request.Body)
			}

      response = Response(nil, CREATED, nil)
		}

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

