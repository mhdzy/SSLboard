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
	"strings"
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
func verifyLogin(sslClient pb.SSLboardClient) string {

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

	// TODO: reset password so that it isn't in memory

	// empty print for formatting
	fmt.Println()

	return username

}

/**
 * func parse
 * Parses a command into separate strings
 */
func parse(cmd string) (string, string, string, bool) {

	words := strings.Split(cmd, " ")

	command := words[0] // check to make sure this exists
	group := words[1]   // check to make sure this exists
	message := words[2:len(words)]
	err := false

	return command, group, strings.Join(message, " "), err

}

/**
 * func interactWithBoard
 * Authenticates user, takes arguments from command line
 */
func interactWithBoard(username string, sslClient pb.SSLboardClient) {

	flg := true
	reader := bufio.NewReader(os.Stdin)

	for flg {

		// prompt
		fmt.Printf("> ")

		// read command ("GET/POST/END GROUP MESSAGE")
		cmd, _ := reader.ReadString('\n')

		// parse cmd into three separate strings
		command, group, message, err := parse(cmd)

		if err == true {
			fmt.Println("rpc syntax: <GET/POST/END> <groupName> <POST: messageContent>")
			continue
		}

		// create struct to send over TLS pipe
		packet := &pb.Message{Username: username, Group: group, Msg: message}
		packet2 := &pb.Credentials{Username: username}

		// GET rpc call
		if command == "GET" {
			sslClient.Get(context.Background(), packet)
		} else if command == "POST" {
			sslClient.Post(context.Background(), packet)
		} else if command == "END" {
			sslClient.End(context.Background(), packet2)
		} else {
			fmt.Println("You issued an incorrect command. <GET/POST/END> are acceptable.")
		}

	}

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
	username := verifyLogin(sslClient)

	// SECOND: allow interaction with the message board
	interactWithBoard(username, sslClient)

	fmt.Println("Exiting client.")

}
