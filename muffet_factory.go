package main

type muffetOptions struct {
	arguments []string
}

type muffetFactory interface {
	Create(options muffetOptions) muffetExecutor
}
