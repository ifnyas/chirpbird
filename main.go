package main

import (
	"chirpbird/chat"
	"chirpbird/router"
)

func main() {
	go chat.Init(2)
	router.Run()
}
