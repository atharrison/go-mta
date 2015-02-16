package main

import (
	"fmt"
	"runtime"
)


func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Starting the Go MTA Server.\n")

	// Start the Server listener
	go server()

	// Die after input is read.
	var input string
	fmt.Scanln(&input)
}
