package server

import (
	"fmt"
	"io"
	"net"
	"os"
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
		request := make([]byte, 1024)
		reqBytes, err := conn.Read(request)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading data from the conn: ", err.Error())
			return
		}

		fmt.Printf("Reques received with %d bytes\n\n", reqBytes)
		fmt.Println(string(request[:reqBytes]))

    response := []byte("HTTP/1.1 200 OK\r\n\r\n")
		bytes, err := conn.Write(response)
		if err != nil {
			fmt.Println("Error sending data to the connection: ", err.Error())
			return
		}

		fmt.Printf("Send %v bytes\n", bytes)
	}
}
