//
//  client.go
//  CS 419
//

package main

import (
	"crypto/tls"
	"log"
	"net"
	"os"
)

var STATUS bool = true

/**
 * func connectToBoard
 * connects to the message board via net.Dial
 */
func connectToServer(addr string) *tls.Conn {

	// initialize config struct (ref)
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	// initiate server connection
	conn, err := tls.Dial("tcp", addr, config)
	if err != nil {
		log.Println(err)
		panic("Error in tls.Dial().")
	}

	log.Println("Client successfully connected to server via TLS.")

	return conn
}

/**
 * func interactWithBoard
 * authenticates user, takes arguments from command line
 */
func interactWithBoard(conn net.Conn) {

	defer conn.Close()

	// get stuff from command line (user, pass)

	// bundle into JSON

	// write over connection

	// read ack

	// forloop

	// get commandline input

	// put input into JSON

	// write command

	// read response

	// print response ?

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

	interactWithBoard(conn)

	log.Print(conn)

}
