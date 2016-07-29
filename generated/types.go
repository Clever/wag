package main

type Book struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

func (v Book) Validate() error { return nil }
type Error struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

func (v Error) Validate() error { return nil }
