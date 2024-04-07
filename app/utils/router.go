package utils

import (
	"fmt"
	"regexp"
)

type Handler func(request *Request) []byte

type Router struct {
	getMap          map[string]Handler
	postMap         map[string]Handler
	NotFoundHandler Handler
}

func (r *Router) Get(path string, handler Handler) {
	if r.getMap == nil {
		r.getMap = make(map[string]Handler)
	}

	r.getMap[path] = handler
}

func (r *Router) Post(path string, handler Handler) {
	if r.postMap == nil {
		r.postMap = make(map[string]Handler)
	}

	r.postMap[path] = handler
}

func (r *Router) Exec(request *Request) []byte {
	var executingMap map[string]Handler

	switch request.Method {
	case GET:
		executingMap = r.getMap
	case POST:
		executingMap = r.postMap
	}

	fmt.Println(executingMap)
	for matcher, handler := range executingMap {
		itMatches := regexp.MustCompile(matcher).MatchString(request.Path)

		if !itMatches {
			continue
		}

		fmt.Println("path matcher ", matcher, handler)

		return handler(request)
	}

	return r.NotFoundHandler(request)
}
