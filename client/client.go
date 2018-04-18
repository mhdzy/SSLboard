//
//  client.go
//  CS 419
//

package main

import (
	"os"
//	"fmt"
//	"net"
	"log"
	"crypto/tls"
)

const ADDR string = "127.0.0.1:8080" // ADDRESS:PORT
var STATUS bool = true


func connectToServer(addr string) {

	// begin with empty config struct (ref)
	config := &tls.Config{}

	// initiate server connection
	conn, err := tls.Dial("tcp", addr, config)
	// conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
        return
	}
	defer conn.Close()
	
	log.Println("Client successfully connected to server.")

	// do things

	// close connection
	log.Println("Client closing connection.")
}

// should connect to server and transfer messages
func main() {

	if len(os.Args) != 2 {
        log.Fatal(os.Stderr, "Usage: %s <ip-addr>\n", os.Args[0])
        os.Exit(1)
    }

    addr := os.Args[1]
    log.Println("Connecting to ", addr)

    connectToServer(addr)

}
