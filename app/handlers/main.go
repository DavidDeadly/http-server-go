package handlers

import (
	"fmt"

	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

const (
	OK          = 200
	CREATED     = 201
	NOT_FOUND   = 404
	BAD_REQUEST = 400
)

func Response(body *string, code uint16, contentType *string) []byte {
	response := fmt.Sprintf("HTTP/1.1 %d OK\r\n", code)

	if contentType == nil {
		content := "text/plain"
		contentType = &content
	}

	if body != nil {
		response += fmt.Sprintf("Content-Type: %s\r\n", *contentType)
		response += fmt.Sprintf("Content-Length: %d\r\n", len(*body))
		response += fmt.Sprintf("\r\n%s", *body)
	}

	response += "\r\n"

	fmt.Println(response)
	return []byte(response)
}

func NotFoundResponse(request *utils.Request) []byte {
	response := fmt.Sprintf("HTTP/1.1 %d Not Found\r\n\r\n", NOT_FOUND)

	return []byte(response)
}
