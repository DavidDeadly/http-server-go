package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/config"
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	config.Config()

	listener := server.StartServerOn("4221")

	server.ListenForConnections(listener)
}
