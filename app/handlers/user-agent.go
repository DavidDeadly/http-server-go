package handlers

import "github.com/codecrafters-io/http-server-starter-go/app/utils"

func UserAgent(request *utils.Request) []byte {
	userAgent, ok := request.Headers["User-Agent"]

	if !ok {
		return Response(nil, BAD_REQUEST, nil)
	}

	return Response(&userAgent, OK, nil)
}
