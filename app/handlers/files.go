package handlers

import (
	"fmt"
	"os"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

func Files(request *utils.Request) []byte {
	fileName := strings.Replace(request.Path, "/files/", "", 1)
	filePath := fmt.Sprintf("%s/%s", config.CONFIG[config.DIR], fileName)

	switch {
	case request.Method == GET:
		return readFile(filePath)
	case request.Method == POST:
		return writeFile(filePath, request.Body)
	}

	return []byte{}
}

func readFile(filePath string) []byte {
	fmt.Println("Reading... ", filePath)

	data, err := os.ReadFile(filePath)

	if err == nil {
		body := string(data[0:])
		contentType := "application/octet-stream"
		return Response(&body, OK, &contentType)
	}

	return NotFoundResponse()
}

func writeFile(filePath string, content string) []byte {
	err := os.WriteFile(filePath, []byte(content), 0644)
	// TODO: add support for create files inside directories
	fmt.Println("Writing ... ", filePath, " with: ", content)

	if err != nil {
		fmt.Printf("Error writing the file with the content: %s\n", content)
	}

	return Response(nil, CREATED, nil)
}
