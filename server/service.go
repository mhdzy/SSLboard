package main

import (
	"context"
	"fmt"
	"log"

	"github.com/SleightOfHandzy/SSLboard/pb"
	"golang.org/x/crypto/bcrypt"
)

type SSLboardServer struct{}

/**
 * func Authenticate
 * Authenticates a given username/passwords
 */
func (s *SSLboardServer) Authenticate(ctx context.Context, c *pb.Credentials) (*pb.Credentials, error) {

	log.Println("RPC call to service.Authenticate")

	// extract username from the struct passed through TLS
	usr := c.Username
	pwd := []byte(c.Password)

	// add salt and hash
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}

	// print out username and password and hash for debugging purposes
	fmt.Printf("*debug* username: %s", usr)
	fmt.Printf("*debug* password: %s\n", string(pwd))
	fmt.Printf("*debug* hash pwd: %s\n", string(hash))

	// compare stored hashed password and password from database
	err = bcrypt.CompareHashAndPassword(hash, pwd)
	if err != nil {
		log.Println(err)
	}

	// \n is used to create formatting
	log.Printf("Authenticated user: %s\n", usr)

	// PLEASE SEND BACK UNIQUE TOKEN ID FOR USER

	// compare with stored db credentials

	return c, nil
}

/**
 * func get
 * Handles a GET request from the client
 */
func (s *SSLboardServer) Get(_ context.Context, m *pb.Message) (*pb.Message, error) {

	log.Println("RPC call to service.Get")
	log.Printf("Username: %s", m.Username) // m.Username CONTAINS A \n
	log.Printf("Group: %s", m.Group)       // m.Group CONTAINS A \n
	log.Printf("Message: %s\n", m.Msg)     // \n for formatting

	return m, nil

}

/**
 * func post
 * Handles a POST request from the client
 */
func (s *SSLboardServer) Post(_ context.Context, m *pb.Message) (*pb.Message, error) {

	log.Println("RPC call to service.Post")
	log.Printf("Username: %s", m.Username) // m.Username CONTAINS A \n
	log.Printf("Group: %s", m.Group)       // m.Group CONTAINS A \n
	log.Printf("Message: %s", m.Msg)       // m.Msg CONTAINS A \n

	return m, nil

}

/**
 * func end
 * Handles a GET request from the client
 */
func (s *SSLboardServer) End(_ context.Context, c *pb.Credentials) (*pb.Credentials, error) {

	log.Println("RPC call to service.End")
	log.Printf("Username: %s", c.Username) // m.Username CONTAINS A \n

	// remove client token from list of active tokens

	return c, nil

}
