package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SleightOfHandzy/SSLboard/pb"
	"github.com/boltdb/bolt"
	"golang.org/x/crypto/bcrypt"
)

type SSLboardServer struct{}

/**
 * func Authenticate
 * Authenticates a given username/passwords
 */
func (s *SSLboardServer) Authenticate(ctx context.Context, c *pb.Credentials) (*pb.Credentials, error) {

	log.Println("RPC call to service.Authenticate")

	// open database
	db, err := bolt.Open("./board.db", 0777, nil)
	if err != nil {
		return c, err
	}
	defer db.Close()

	// extract username from the struct passed through TLS
	var hash []byte
	bucket_name := []byte("Users")
	username := []byte(c.Username)
	password := []byte(c.Password)

	// get username from database
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket_name)
		if err != nil {
			log.Fatal(err)
		}

		// get hash and salt from user bucket
		stored_hash := bucket.Get(username)
		if stored_hash == nil {
			log.Println("Username does not exist; adding now")
			return errors.New("Username does not exist")
		}
		hash = make([]byte, len(stored_hash))
		n := copy(hash, stored_hash)
		fmt.Println(n)

		return nil
	})

	// username does not exist: create new key pair
	if err != nil {
		hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(hash))
		err = bcrypt.CompareHashAndPassword(hash, password)
		if err != nil {
			log.Fatal(err)
		}

		// store username and hash
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucket_name)
			if err != nil {
				log.Fatal(err)
			}

			err = bucket.Put(username, hash)
			if err != nil {
				return err
			}
			hash = bucket.Get(username)
			fmt.Println(string(hash))

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Added new user")
		return c, err
	} else {
		// compare stored hashed password and password from database
		fmt.Println(string(hash))
		err = bcrypt.CompareHashAndPassword(hash, password)
		if err != nil {
			// username exists but password doesn't match: return error
			log.Println("Username exists but password doesn't match: return error")
			return c, err
		} else {
			// user is authenticated: send back unique token
			log.Println("User is authenticated: send back unique token")
			return c, err
		}
	}

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
