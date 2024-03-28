package handlers

import (
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/utils"
)

func Echo(request *utils.Request) []byte {
	string := strings.Replace(request.Path, "/echo/", "", 1)

	return Response(&string, OK, nil)
}
