//
//  server.go
//  CS 419
//

package main

import (
	"crypto/tls"
	"log"
	"net"
)

const PORT string = ":8080" // PORT
var STATUS bool = true      // status of server

/**
 * func main
 * Listens for connections, accepts, and hands off to acceptor
 */
func main() {

	log.Printf("Server running on %s", PORT)

	// load certificate from files
	cert, err := tls.LoadX509KeyPair("../server.crt", "../server.key")
	if err != nil {
		log.Println(err)
		panic("Error in loading x509 key pairs.")
	}

	// create config struct (ref) from certificate
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// listens for connections on port
	ln, err := tls.Listen("tcp", PORT, config)
	if err != nil {
		log.Println(err)
		panic("Error in listening on port.")
	}
	defer ln.Close()

	// loop while server (status) is true
	for STATUS {

		// accept an incoming connection
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
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

	defer conn.Close()

	// print locally that a client is connected
	log.Println("Accepted a client connection.")

	// transmit an acceptance message to the client
	conn.Write([]byte("Accepted\n"))

	// do some work
	for {

	}

	// close connection
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
