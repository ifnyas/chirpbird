package main

import (
	"chirpbird/utils"
)

func main() {
	go utils.H.Run()
	utils.RouterRun()
}
