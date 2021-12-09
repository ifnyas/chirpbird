package main

import (
	"chirpbird/router"
	"chirpbird/wss"
)

func main() {
	go wss.Init()
	router.Run()
}
