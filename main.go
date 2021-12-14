package main

import (
	"chirpbird/router"
	"chirpbird/wss"
)

func main() {
	go wss.Init(0)
	router.Run()
}
