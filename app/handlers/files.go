package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

func ReadFile(request *utils.Request) []byte {
	fileName := strings.Replace(request.Path, "/files/", "", 1)
	filePath := fmt.Sprintf("%s/%s", config.CONFIG[config.DIR], fileName)
	fmt.Println("Reading... ", filePath)

	data, err := os.ReadFile(filePath)

	if err == nil {
		body := string(data[0:])
		contentType := "application/octet-stream"
		return Response(&body, OK, &contentType)
	}

	return NotFoundResponse(nil)
}

func WriteFile(request *utils.Request) []byte {
	fileName := strings.Replace(request.Path, "/files/", "", 1)
	filePath := fmt.Sprintf("%s/%s", config.CONFIG[config.DIR], fileName)

	err := os.WriteFile(filePath, []byte(request.Body), 0644)
	// TODO: add support for create files inside directories
	fmt.Println("Writing ... ", filePath, " with: ", request.Body)

	if err != nil {
		fmt.Printf("Error writing the file with the content: %s\n", request.Body)
	}

	return Response(nil, CREATED, nil)
}
