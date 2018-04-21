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

	log.Println("Authenticating user...")

	usr := c.Username
	pwd := []byte(c.Password)

	// add salt and hash
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", usr)
	fmt.Printf("%s\n", string(pwd))
	fmt.Printf("%s\n", string(hash))

	err = bcrypt.CompareHashAndPassword(hash, pwd)
	if err != nil {
		log.Println(err)
	}
	log.Println("Authenticates")

	// compare with stored db credentials

	return c, nil
}

/**
 * func get
 * Handles a GET request from the client
 */
func (s *SSLboardServer) Get(_ context.Context, m *pb.Message) (*pb.Message, error) {

	return m, nil

}

/**
 * func post
 * Handles a POST request from the client
 */
func (s *SSLboardServer) Post(_ context.Context, m *pb.Message) (*pb.Message, error) {

	return m, nil

}

/**
 * func end
 * Handles a GET request from the client
 */
func (s *SSLboardServer) End(_ context.Context, c *pb.Credentials) (*pb.Credentials, error) {

	return c, nil

}
