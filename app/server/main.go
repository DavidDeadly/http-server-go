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

const (
  OK = 200
  NOT_FOUND = 404
  BAD_REQUEST = 400
)

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

			response = Response(nil, OK)

		case regexp.MustCompile("/echo/*").MatchString(request.path):
			string := strings.Replace(request.path, "/echo/", "", 1)

			response = Response(&string, OK)

		case request.path == "/user-agent":
      userAgent , ok := request.headers["User-Agent"]

      if !ok {
        response = Response(nil, BAD_REQUEST)
      }

			response = Response(&userAgent, OK)

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
	headers map[string]string
	version string
}

func Response(body *string, code uint16) []byte {
	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n", code)

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
	response := fmt.Sprintf("HTTP/1.1 %d Not Found\r\n\r\n", NOT_FOUND)

	return []byte(response)
}

func ParseHTTPRequest(request []byte, requestBytes int) Request {
	stringRequest := string(request[:requestBytes])

	fmt.Printf("Request received with %d bytes\n\n", requestBytes)
	fmt.Println(stringRequest)

  var method string
  var headers = map[string]string{}

	lines := strings.Split(stringRequest, "\r\n")

	startLine := strings.Fields(lines[0])
	method, path, version := startLine[0], startLine[1], startLine[2]

  for _, header := range lines[1:] {
    if header == "" { break }

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