package main

import (
	"fmt"
	"net"
)

const ADDR string = "127.0.0.1:8080" // PORT
var STATUS bool = true

// should connect to server and transfer messages
func main() {

	// startup messages
	fmt.Println("Client booting...")
	fmt.Println("Client connecting to server...")

	// initiate server connect
	_, err := net.Dial("tcp", ADDR)
	if err != nil {
		// handle error
	}

	// connect success!
	fmt.Println("Client successfully connected to server.")

	// accept commands and pass to server
	for {

	}

}
