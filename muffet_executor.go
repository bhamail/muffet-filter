package main

type muffetExecutor interface {
	Check() (string, error)
}
