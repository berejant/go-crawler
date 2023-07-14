package main

type ParserInterface interface {
	Init() error
	ProcessPage(url string) (*ParserResponse, error)
}

type ParserResponse struct {
	Links []string
}
