package main

type Result struct {
	Address string
	Error   error
}

func (r Result) isOk() bool {
	return r.Error == nil
}
