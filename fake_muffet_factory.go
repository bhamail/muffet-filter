package main

type fakeMuffetFactory struct {
	response fakeMuffetResponse
}

func newFakeMuffetFactory(message string, err error) *fakeMuffetFactory {
	resp := fakeMuffetResponse{message, err}
	return &fakeMuffetFactory{response: resp}
}
func (f *fakeMuffetFactory) Create(options muffetOptions) muffetExecutor {
	return &fakeMuffetExecutor{options}
}

type fakeMuffetExecutor struct {
	options muffetOptions
}

//goland:noinspection GoUnusedParameter
func (r *fakeMuffetExecutor) Check(args *arguments) (string, error) {
	return "[]", nil
}
