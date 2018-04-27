package main

import (
	"context"
	"errors"
	"log"
	"math/rand"

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
		copy(hash, stored_hash)

		return nil
	})

	if err != nil {

		// create new key pair
		hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
		if err != nil {
			log.Fatal(err)
		}

		// store username and hash
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucket_name)
			if err != nil {
				return err
			}

			err = bucket.Put(username, hash)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Added new user")

	} else {

		// check if user is currently in a session
		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("Tokens"))
			if bucket == nil {
				return errors.New("Bucket does not exist")
			}
			token := bucket.Get(username)
			if token != nil {
				return errors.New("User is currently in a session")
			}
			return nil
		})

		if err != nil {
			log.Fatal(err)
			return c, err
		}

		// compare stored hashed password and password from database
		err = bcrypt.CompareHashAndPassword(hash, password)

		if err != nil {
			log.Println("Incorrect password")
			return c, err
		}
	}

	log.Println("User is authenticated")

	// generate token
	token := make([]byte, 16)
	n, err := rand.Read(token)
	if n != 16 {
		log.Fatal(err)
	}

	// store token
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Tokens"))
		if err != nil {
			log.Fatal(err)
		}

		err = bucket.Put(username, token)
		if err != nil {
			return err
		}

		return nil
	})

	// return token
	c.Password = string(token)
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
