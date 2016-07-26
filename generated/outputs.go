package main

type GetBookOutput interface {
	GetBookStatus() int
}

type GetBook200Output Book

func (o GetBook200Output) GetBookStatus() int {
	return 200
}

type GetBookDefaultOutput Error

func (o GetBookDefaultOutput) GetBookStatus() int {
	return 200
}

