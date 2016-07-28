package main

type Book struct {
	Id string `json:"id"`
	Name string `json:"name"`
}

type Error struct {
	Message string `json:"message"`
	Code int `json:"code"`
}

