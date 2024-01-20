package main

type muffetOptions struct {
}

type muffetFactory interface {
	Create(options muffetOptions) muffetExecutor
}
