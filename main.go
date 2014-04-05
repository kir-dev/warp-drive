package main

import (
	"fmt"
)

var (
	config configuration
)

func main() {
	println("Hello world!")

	// TODO: make path configrable by command line args
	config = loadConfiguration("config/config.json")
	fmt.Println(config)
}
