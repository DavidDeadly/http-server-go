package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	listener := server.StartServerOn("4221")

	server.ListenForConnections(listener)
}
