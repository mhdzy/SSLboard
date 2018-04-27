//
//  server.go
//  CS 419
//

package main

import (
	"crypto/tls"
	"log"
	"net"

	"github.com/SleightOfHandzy/SSLboard/pb"
	"google.golang.org/grpc"
)

const PORT string = ":8080" // PORT
var STATUS bool = true      // status of server

/**
 * func loadKeyPairs()
 * Loads cert and key pairs.
 */
func loadKeyPairs() tls.Certificate {
	cert, err := tls.LoadX509KeyPair("./server.crt", "./server.key")

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
 * func main
 * Listens for connections, accepts, and hands off to acceptor
 */
func main() {

	log.Printf("Server running on %s.\n", PORT)

	// load Bolt db
	// db, err := bolt.Open("./board.db", 0777, nil)
	// if err != nil {
	// 	panic("Database failed to open")
	// }
	// defer db.Close()
	//
	// db.Update(func(tx *bolt.Tx) error {
	// 	tx.DeleteBucket([]byte("Tokens"))
	// 	return nil
	// })

	// load certificate from files
	cert := loadKeyPairs()

	// create config struct (ref) from certificate
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	// listen on the given port
	ln := listener(config)

	srv := grpc.NewServer()
	pb.RegisterSSLboardServer(srv, &SSLboardServer{})

	srv.Serve(ln)

}
