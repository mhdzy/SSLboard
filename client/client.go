//
//  client.go
//  CS 419
//

package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"syscall"

	"github.com/SleightOfHandzy/SSLboard/pb"
	"golang.org/x/crypto/ssh/terminal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var STATUS bool = true

/**
 * func connectToServer
 * Connects to the message board via net.Dial
 */
func connectToServer(addr string) (pb.SSLboardClient, *grpc.ClientConn) {

	// initialize config struct (ref)
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	creds := credentials.NewTLS(config)

	// initiate server connection
	grpcConn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Println(err)
		panic("Error in grpc.Dial().")
	}

	log.Println("Client successfully connected to server via TLS.")

	sslClient := pb.NewSSLboardClient(grpcConn)

	return sslClient, grpcConn
}

/**
 * func verifyLogin
 * Verifies a username/password combination in the server's database.
 */
func verifyLogin(sslClient pb.SSLboardClient) {

	// establish a reader to read username
	reader := bufio.NewReader(os.Stdin)

	// get username from command line
	fmt.Printf("username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		fmt.Println("Error in reading username, setting to default 'pk419'.")
	}

	// get password securely from command line
	fmt.Printf("password: ")
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Println(err)
		fmt.Println("Error in reading password, setting to default 'spring18'")
	}

	// create credentials to pass through TLS pipe in rpc call
	cred := &pb.Credentials{Username: username, Password: string(password)}

	// RPC authentication method (send credentials)
	_, err = sslClient.Authenticate(context.Background(), cred)
	if err != nil {
		log.Println(err)
		panic("Error in sslClient.Authenticate rpc call.")
	}

	// empty print for formatting
	fmt.Println()

}

/**
 * func interactWithBoard
 * Authenticates user, takes arguments from command line
 */
func interactWithBoard(sslClient pb.SSLboardClient) {

	// forloop

	// take input from command line

	// set pb struct

	// do RPC call

}

/**
 * func main
 * Orchestrates the client's interaction with the server
 */
func main() {

	// checks that a IP address was specified
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <ip-addr>.\n", os.Args[0])
		panic("*E* Error in command line args.")
	}

	addr := os.Args[1]
	log.Println("Connecting to: ", addr)

	// conn is of type pb.SSLboardClient
	sslClient, grpcConn := connectToServer(addr)
	defer grpcConn.Close()

	// FIRST: verify username/password combination
	verifyLogin(sslClient)

	// SECOND: allow interaction with the message board
	interactWithBoard(sslClient)

	fmt.Println("Exiting client.")

}
