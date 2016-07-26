package main

type Book struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type Error struct {
	Code int `json:"code"`
	Message string `json:"message"`
}

