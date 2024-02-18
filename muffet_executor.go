package main

type muffetExecutor interface {
	Check(args *arguments) (string, error)
}
