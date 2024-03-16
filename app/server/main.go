package server

import (
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
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

func handleConnection(conn net.Conn) {
	fmt.Println("Handling...")
	defer conn.Close()

	for {
		rawRequest := make([]byte, 1024)
		bytes, err := conn.Read(rawRequest)

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading data from the conn: ", err.Error())
			return
		}

		request := ParseHTTPRequest(rawRequest, bytes)

		var response []byte

    

    switch {

    case request.path == "/":
      response = SuccessResponse(nil)

    case regexp.MustCompile("/echo/*").MatchString(request.path):
      string := strings.Replace(request.path, "/echo/", "", 1)

      response = SuccessResponse(&string)

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
}

type METHOD string

const (
	GET  METHOD = "GET"
	POST METHOD = "POST"
)

type Request struct {
	method  METHOD
	path    string
	version string
}

func SuccessResponse(body *string) []byte {
  response := "HTTP/1.1 200 OK\r\n"

  if body != nil {
    response += "Content-Type: text/plain\r\n"
    response += fmt.Sprintf("Content-Length: %d\r\n", len(*body))
    response += fmt.Sprintf("\r\n%s", *body)
  }

  response += "\r\n"

  fmt.Println(response)
	return []byte(response)
}

func NotFoundResponse() []byte {
	return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
}

func ParseHTTPRequest(request []byte, requestBytes int) Request {
	stringRequest := string(request[:requestBytes])

	fmt.Printf("Request received with %d bytes\n\n", requestBytes)
	fmt.Println(stringRequest)

	lines := strings.Split(stringRequest, "\r\n")

	startLine := strings.Fields(lines[0])
	method, path, version := startLine[0], startLine[1], startLine[2]

	return Request{
		method:  METHOD(method),
		path:    path,
		version: version,
	}
}
