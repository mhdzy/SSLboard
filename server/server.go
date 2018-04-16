//
//  client.go
//  CS 419
//
//  Created by Matthew Handzy on 4/16/18.
//  Copyright © 2018 Matthew Handzy. All rights reserved.
//

package main

import (
	"fmt"
	"net"
)

const PORT string = ":8080" // PORT
var STATUS bool = true      // status of server

/**
 * func main
 * Listens for connections, accepts, and hands off to acceptor
 */
func main() {

	fmt.Println("Server booting...") // message to ensure the program started

	// listens for connections on port
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		// handle error
	}

	// loop while server (status) is true
	for STATUS {

		// accept an incoming connection
		conn, err := ln.Accept()
		if err != nil {
			//handle error
		}

		// pass the Connection object to our acceptor
		// this is a thread call
		go accept(conn)
	}

}

/**
 * func accept
 * accepts a Connection object to handle the TCP messages
 */
func accept(conn net.Conn) {

	// print locally that a client is connected
	fmt.Println("Accepted a client connection.")

	// transmit an acceptance message to the client
	//conn.Write([]byte("Accepted.\n"))

	// then go do some work
}

// handles the GET command
func get(group string) {
}

// handles the POST command
func post(message string, group string) {
}

// handles the END command
func end() {
}
