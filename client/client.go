//
//  client.go
//  CS 419
//

package main

import (
	"crypto/tls"
	"log"
	"os"
)

var STATUS bool = true

/**
 * func connectToBoard
 * connects to the message board via net.Dial
 */

func connectToServer(addr string) *tls.Conn {

	// begin with empty config struct (ref)
	config := &tls.Config{}

	// initiate server connection
	conn, err := tls.Dial("tcp", addr, config)

	// gets stuck on tls.Dial

	// conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		panic("Error in tls.Dial().")
	}

	log.Println("Client successfully connected to server via TLS.")

	return conn
}

// should connect to server and transfer messages
func main() {

	// checks that a IP address was specified
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <ip-addr>\n", os.Args[0])
		panic("*E* Error in command line args.")
	}

	addr := os.Args[1]
	log.Println("Connecting to: ", addr)

	conn := connectToServer(addr)

	log.Print(conn)

}
