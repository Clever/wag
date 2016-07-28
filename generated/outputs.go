package main

type GetBookByIDOutput interface {
	GetBookByIDStatus() int
}

type GetBookByID200Output Book

func (o GetBookByID200Output) GetBookByIDStatus() int {
	return 200
}

type GetBookByIDDefaultOutput Error

func (o GetBookByIDDefaultOutput) GetBookByIDStatus() int {
	return 200
}

