package cqrs

type Request struct {
	Service string
	Head    any
	Body    any
}
