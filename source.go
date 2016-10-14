package main

type Object struct {
	Key string
}

type Source interface {
	Run()
}
