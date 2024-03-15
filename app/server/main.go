package server

import (
	"fmt"
	"io"
	"net"
	"os"
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

		if request.path == "/" {
      response = SuccessResponse()
		} else {
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

func SuccessResponse() []byte {
	return []byte("HTTP/1.1 200 OK\r\n\r\n")
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
