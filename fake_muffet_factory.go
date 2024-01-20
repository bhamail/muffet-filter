package main

type fakeMuffetFactory struct {
	response fakeMuffetResponse
}

func newFakeMuffetFactory(message string, err error) *fakeMuffetFactory {
	resp := fakeMuffetResponse{message, err}
	return &fakeMuffetFactory{response: resp}
}
func (f *fakeMuffetFactory) Create(options muffetOptions) muffetExecutor {
	return f.Create(options)
}
