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
 * func loadKeyPairs()
 * Loads cert and key pairs.
 */
func loadKeyPairs() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("../server.crt", "../server.key")

	if err != nil {
		log.Println(err)
		panic("Error in loading x509 key pairs.")
	}

	return cert
}

/**
 * func listener
 * Creates a listener over TLS and returns that secure object
 */
func listener(config *tls.Config) net.Listener {
	ln, err := tls.Listen("tcp", PORT, config)
	if err != nil {
		log.Println(err)
		panic("Error in listening on port.")
	}
	return ln
}

/**
 * func get
 */
func get(group string) {
}

/**
 * func post
 */
func post(message string, group string) {
}

/**
 * func end
 */
func end() {
}

/**
 * func connectionHandler
 * handles all connection work!
 */
func connectionHandler(conn net.Conn) {
	defer conn.Close()

	// print locally that a client is connected
	log.Println("Accepted a client connection.")

	// transmit an acceptance message to the client
	conn.Write([]byte("Accepted.\n"))

	// do some work
	//for { /* COMMENTED OUT FOR LOOP SINCE SERVER WOULD HANG */ }

}

/**
 * func main
 * Listens for connections, accepts, and hands off to acceptor
 */
func main() {

	log.Printf("Server running on %s.", PORT)

	// load certificate from files
	cert := loadKeyPairs()

	// create config struct (ref) from certificate
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// listen on the given port
	ln := listener(config)

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
		go connectionHandler(conn)
	}

}
