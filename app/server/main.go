package server

import (
	"fmt"
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
		_, err := listener.Accept()
		if err != nil {
			fmt.Println(err, "Error accepting connection")
			os.Exit(1)
		}

		fmt.Println("\nNew connection")

		// go handleConnection(connection)
	}
}
