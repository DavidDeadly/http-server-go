package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
)

func StartServerOn(port string) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Server running on PORT: %s\n", port)

	return listener
}

func ListenForConnections(listener net.Listener) net.Conn {
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err, "Error accepting connection")
			os.Exit(1)
		}

		fmt.Println("\nNew connection")

		go handleConnection(connection)
	}
}

const (
	OK          = 200
	NOT_FOUND   = 404
	BAD_REQUEST = 400
)

func handleConnection(conn net.Conn) {
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

	request := ParseHTTPRequest(rawRequest, bytes)

	var response []byte

	switch {

	case request.path == "/":

		response = Response(nil, OK, nil)

	case regexp.MustCompile("/echo/*").MatchString(request.path):
		string := strings.Replace(request.path, "/echo/", "", 1)

		response = Response(&string, OK, nil)

	case request.path == "/user-agent":
		userAgent, ok := request.headers["User-Agent"]

		if !ok {
			response = Response(nil, BAD_REQUEST, nil)
		} else {
      response = Response(&userAgent, OK, nil)
    }

	case regexp.MustCompile("/files/*").MatchString(request.path):
		fileName := strings.Replace(request.path, "/files/", "", 1)

		filePath := fmt.Sprintf("%s/%s", config.CONFIG[config.DIR], fileName)

		fmt.Println("Reading... ", filePath)
		data, err := os.ReadFile(filePath)
		if err == nil {
      body := string(data[0:])
      contentType := "application/octet-stream"
      response = Response(&body, OK, &contentType)
		} else {
			response = NotFoundResponse()
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

type METHOD string

const (
	GET  METHOD = "GET"
	POST METHOD = "POST"
)

type Request struct {
	method  METHOD
	path    string
	headers map[string]string
	version string
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

func ParseHTTPRequest(request []byte, requestBytes int) Request {
	stringRequest := string(request[:requestBytes])

	fmt.Printf("Request received with %d bytes\n\n", requestBytes)
	fmt.Println(stringRequest)

	var method string
	headers := map[string]string{}

	lines := strings.Split(stringRequest, "\r\n")

	startLine := strings.Fields(lines[0])
	method, path, version := startLine[0], startLine[1], startLine[2]

	for _, header := range lines[1:] {
		if header == "" {
			break
		}

		data := strings.Split(header, ": ")
		key, value := data[0], data[1]

		headers[key] = value
	}

	return Request{
		method:  METHOD(method),
		path:    path,
		version: version,
		headers: headers,
	}
}
