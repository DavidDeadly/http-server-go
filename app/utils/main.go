package utils

import (
	"fmt"
	"strings"
)

type METHOD string

const (
	GET  METHOD = "GET"
	POST METHOD = "POST"
)

type Request struct {
	Method  METHOD
	Path    string
	Headers map[string]string
	Version string
	Body    string
}

func ParseHTTPRequest(request []byte, requestBytes int) *Request {
	stringRequest := string(request[:requestBytes])

	fmt.Printf("Request received with %d bytes\n\n", requestBytes)
	fmt.Println(stringRequest)

	var method string
	headers := map[string]string{}

	lines := strings.Split(stringRequest, "\r\n")

	startLine := strings.Fields(lines[0])
	method, path, version := startLine[0], startLine[1], startLine[2]

	var i int
	for _, header := range lines[1:] {
		if header == "" {
			break
		}

		data := strings.Split(header, ": ")
		key, value := data[0], data[1]

		headers[key] = value
		i++
	}

	body := strings.Join(lines[i+2:], "\r\n")

	return &Request{
		Method:  METHOD(method),
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}
}
