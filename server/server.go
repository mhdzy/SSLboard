//
//  server.go
//  CS 419
//

package main

import (
//	"fmt"
	"net"
	"log"
	"crypto/tls"
)

const PORT string = ":8080" // PORT
var STATUS bool = true      // status of server

/**
 * func main
 * Listens for connections, accepts, and hands off to acceptor
 */
func main() {

	log.Println("Server running on %s", PORT)

	cert, err := tls.LoadX509KeyPair("../server.crt", "../server.key")
    if err != nil {
        log.Println(err)
        return
    }
    log.Println(cert)

	// listens for connections on port
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Println(err)
        return
	}

	// loop while server (status) is true
	for STATUS {

		// accept an incoming connection
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
        	return
		}

		// pass the Connection object to our acceptor
		// this is a thread call
		go clientThread(conn)
	}

}

/**
 * func clientThread
 */
func clientThread(conn net.Conn) {

	// print locally that a client is connected
	log.Println("Accepted a client connection.")

	// transmit an acceptance message to the client
	//conn.Write([]byte("Accepted.\n"))

	// do some work

	// close connection
	conn.Close()
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
