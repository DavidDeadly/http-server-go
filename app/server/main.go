package server

import (
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/http-server-starter-go/app/handlers"
	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

var router = new(utils.Router)

func StartServerOn(port string) net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		fmt.Println("Failed to bind to port 6379: ", err.Error())
		os.Exit(1)
	}

	router.NotFoundHandler = handlers.NotFoundResponse

	router.Get("^/$", func(request *utils.Request) []byte {
    return handlers.Response(nil, handlers.OK, nil)
  })
	router.Get("/echo/*", handlers.Echo)
	router.Get("/user-agent", handlers.UserAgent)
	router.Get("/files", handlers.ReadFile)

	router.Post("files/*", handlers.WriteFile)

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

	response := router.Exec(request)

	bytes, err = conn.Write(response)
	if err != nil {
		fmt.Println("Error sending data to the connection: ", err.Error())
		return
	}

	fmt.Printf("Send %v bytes\n", bytes)
}
