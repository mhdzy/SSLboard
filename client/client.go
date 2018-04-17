//
//  client.go
//  CS 419
//

package main

import (
	"os"
//	"fmt"
	"net"
	"log"
)

const ADDR string = "127.0.0.1:8080" // ADDRESS:PORT
var STATUS bool = true


func connectToBoard(addr string) {

	log.Println("Client connecting to server...")

	// initiate server connection
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		// handle error
	}
	log.Println("Client successfully connected to server.")

	// do things

	// close connection
	conn.Close()
	log.Println("Client closing connection.")
}

// should connect to server and transfer messages
func main() {

	if len(os.Args) != 2 {
        log.Fatal(os.Stderr, "Usage: %s <ip-addr>\n", os.Args[0])
        os.Exit(1)
    }

    addr := os.Args[1]
    log.Println("The address is ", addr)

    connectToBoard(addr)

}
