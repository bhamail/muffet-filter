package main

type fakeMuffetResponse struct {
	message string
	err     error
}

func (f *fakeMuffetResponse) Response() string {
	return f.message
}

func (f *fakeMuffetResponse) Error() error {
	return f.err
}
