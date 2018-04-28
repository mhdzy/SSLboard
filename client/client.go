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

	sslClient := pb.NewSSLboardClient(grpcConn)

	return sslClient, grpcConn
}

/**
 * func verifyLogin
 * Verifies a username/password combination in the server's database.
 */
func verifyLogin(sslClient pb.SSLboardClient) (string, string) {

	var c *pb.Credentials

	// establish a reader to read username
	reader := bufio.NewReader(os.Stdin)

	// get username from command line
	fmt.Printf("username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Println(err)
		fmt.Println("Error in reading username, setting to default 'pk419'.")
	}

	// remove newline character
	if strings.HasSuffix(username, "\n") {
		username = username[:len(username)-1]
	}

	for {

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
		c, err = sslClient.Authenticate(context.Background(), cred)
		if err != nil {
			fmt.Print("\n")
			log.Println(err)

			// IF ERROR == User is currently in a session, FIX THIS

			continue
		}
		break
	}

	// TODO: reset password so that it isn't in memory

	// empty print for formatting
	fmt.Printf("\nAuthenticated user: %s.\n\n", username)

	// print groups out to the client
	if c.Groups == nil {
		fmt.Printf("There are currently no groups in the messageBoard.\n")
	} else {
		fmt.Printf("Possible groups to GET from and POST to include: \n")
		for i := 0; i < len(c.Groups); i++ {
			fmt.Printf("%s ", c.Groups[i])
		}
		fmt.Printf("\n")
	}

	return username, c.Password
}

/**
 * func parse
 * Parses a command into separate strings
 */
func parse(cmd string) (string, string, string, bool) {

	// remove newline character
	if strings.HasSuffix(cmd, "\n") {
		cmd = cmd[:len(cmd)-1]
	}

	// split input into "cmd grp message"
	words := strings.Split(cmd, " ")

	// extract command from split result
	command := words[0]

	// if END issued, words[n > 1] will segfault
	if command == "END" {
		return command, "", "", false
	}

	// extract groupname from split result
	group := ""
	if len(words) > 1 {
		group = words[1]
	} else {
		return "invalid", "", "", true
	}

	// if GET issued, words[2] "message argument" will not exist
	// if POST issued, words[2:len(words)] includes the message
	// else, an invalid command was issued
	if command == "GET" {
		return command, group, "", false
	} else if command == "POST" {

		// if cmd doesn't contain a message, yell at the user
		if len(words) < 3 {
			fmt.Printf("Please include a message when POSTing to the group.\n")
			return "invalid", "", "", true
		}

		// extract message from strings (cat the rest of the split result)
		message := words[2:len(words)]
		if string(strings.Join(message, " ")) == "" {
			fmt.Printf("Please include a message when POSTing to the group.\n")
			return "invalid", "", "", true
		}
		return command, group, strings.Join(message, " "), false
	} else {
		return "invalid", "", "", true
	}

}

/**
 * func interactWithBoard
 * Authenticates user, takes arguments from command line
 */
func interactWithBoard(username string, token string, sslClient pb.SSLboardClient) {

	reader := bufio.NewReader(os.Stdin)

	// if this function ever ends, the client will issue an END request
	defer func() {
		_, err := sslClient.End(context.Background(),
			&pb.Message{Token: token, Username: username, Group: "", Msg: ""})
		if err != nil {
			fmt.Println(err)
		}
	}()

	// loop through client input
	for {

		// prompt
		fmt.Printf("> ")

		// read command ("GET/POST/END GROUP MESSAGE")
		cmd, _ := reader.ReadString('\n')

		// parse cmd into three separate strings
		command, group, message, err := parse(cmd)

		if err == true {
			fmt.Printf("rpc syntax: <GET/POST/END> <groupName> <(POST call): messageContent>\n\n")
			continue
		}

		// create struct to send over TLS pipe
		packet := &pb.Message{Token: token, Username: username, Group: group, Msg: message}

		// RPC calls are made here
		if command == "GET" {
			c, err := sslClient.Get(context.Background(), packet)

			// check for errors
			if err != nil {
				fmt.Println(err)
			}

			// print messages out to the client
			if c.Messages == nil {
				fmt.Printf("This group does not exist.\n\n")
			} else {
				fmt.Printf("Current Messages: \n")
				for i := 0; i < len(c.Messages); i++ {
					fmt.Printf("%s\n", c.Messages[i])
				}
				fmt.Printf("\n")
			}

			// print groups out to the client
			if c.Groups == nil {
				fmt.Printf("There are currently no groups in the messageBoard.\n")
			} else {
				fmt.Printf("Possible groups to GET from and POST to include: \n")
				for i := 0; i < len(c.Groups); i++ {
					fmt.Printf("%s ", c.Groups[i])
				}
				fmt.Printf("\n")
			}

		} else if command == "POST" {
			c, err := sslClient.Post(context.Background(), packet)

			// check for errors
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("Call to Post() succeeded.\n\n")

			// print groups out to the client
			if c.Groups == nil {
				fmt.Printf("There are currently no groups in the messageBoard.\n")
			} else {
				fmt.Printf("Possible groups to GET from and POST to include: \n")
				for i := 0; i < len(c.Groups); i++ {
					fmt.Printf("%s ", c.Groups[i])
				}
				fmt.Printf("\n")
			}

		} else if command == "END" {
			break // will stack defered function call and exit client
		} else {
			fmt.Printf("You issued an incorrect command. <GET/POST/END> are acceptable.\n\n")
		}
	}
}

/**
 * func main
 * Orchestrates the client's interaction with the server
 */
func main() {

	// get server IP address
	if len(os.Args) != 2 {
		log.Printf("Usage: %s <ip-addr>.\n", os.Args[0])
		panic("*E* Error in command line args.")
	}
	addr := os.Args[1]

	// connect to server using grpc over TLS
	sslClient, grpcConn := connectToServer(addr)
	defer grpcConn.Close()

	// verify username/password combination
	username, token := verifyLogin(sslClient)

	// interact with the message board
	interactWithBoard(username, token, sslClient)

	fmt.Println("Exiting client.")

}
