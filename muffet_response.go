package main

type muffetResponse interface {
	Response() string
	Error() error
}
