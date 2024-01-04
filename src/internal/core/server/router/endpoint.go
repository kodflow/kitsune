package router

import (
	"errors"
	"strings"
)

var RESERVED_ENDPOINTS = map[string]struct{}{
	"public":  {},
	"private": {},
	"doc":     {},
	"docs":    {},
}

type EndPoint struct {
	isRoot   bool
	Depth    int
	Endpoint string
	parent   *EndPoint
	subs     map[string]*EndPoint
	options  []string
}

func (a *EndPoint) Head(...Handler) {
	a.options = append(a.options, "HEAD")
}

func (a *EndPoint) Get(...Handler) {
	a.options = append(a.options, "GET")
}

func (a *EndPoint) Post(...Handler) {
	a.options = append(a.options, "POST")
}

func (a *EndPoint) Put(...Handler) {
	a.options = append(a.options, "PUT")
}

func (a *EndPoint) Patch(...Handler) {
	a.options = append(a.options, "PATCH")
}

func (a *EndPoint) Delete(...Handler) {
	a.options = append(a.options, "DELETE")
}

func (a *EndPoint) Sub(e *EndPoint) *EndPoint {
	if e.parent != nil {
		panic(errors.New("endpoint already has a parent endpoint defined: " + e.Endpoint))
	}

	e.parent = a
	a.subs[e.Endpoint] = e
	e.Depth = a.Depth + 1
	e.updateSubsDepth()

	return e
}

func (a *EndPoint) updateSubsDepth() {
	for _, e := range a.subs {
		e.Depth = a.Depth + 1
		e.updateSubsDepth()
	}
}

func (a *EndPoint) isRootEndpoint() bool {
	return a.isRoot
}

func (a *EndPoint) URL() string {
	if a.isRoot {
		return "/"
	}

	if a.parent.isRoot {
		return "/" + a.Endpoint
	}

	return a.parent.URL() + "/" + a.Endpoint
}

func NewRootPoint() *EndPoint {
	return &EndPoint{
		isRoot:   true,
		Endpoint: "",
		subs:     map[string]*EndPoint{},
		options:  []string{},
	}
}

func NewEndPoint(endpoint string) *EndPoint {
	if strings.Contains(endpoint, "/") {
		panic(errors.New("endpoint can't contain '/'"))
	}

	clearEndpoint := strings.TrimSpace(endpoint)
	if _, exists := RESERVED_ENDPOINTS[clearEndpoint]; exists {
		panic(errors.New("endpoint is reserved: " + clearEndpoint))
	}

	return &EndPoint{
		Endpoint: clearEndpoint,
		subs:     map[string]*EndPoint{},
	}
}

func (e *EndPoint) traverse() []*EndPoint {
	var endpoints []*EndPoint = []*EndPoint{e}

	for _, e := range e.subs {
		endpoints = append(endpoints, e.traverse()...)
	}

	return endpoints
}
