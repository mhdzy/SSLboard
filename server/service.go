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

	log.Println("RPC call to Authenticate()")

	var hash []byte
	var bucket_users = []byte("Users")
	var bucket_tokens = []byte("Tokens")
	var userNotExist = errors.New("Username does not exist")
	var userInSession = errors.New("User is currently in a session")

	// open database
	db, err := bolt.Open("./board.db", 0666, nil)
	if err != nil {
		return c, err
	}
	defer db.Close()

	// extract username from the struct passed through TLS
	username := []byte(c.Username)
	password := []byte(c.Password)

	// get username from database
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket_users)
		if err != nil {
			return err
		}

		// get hash and salt from user bucket
		stored_hash := bucket.Get(username)
		if stored_hash == nil {
			log.Println("Username does not exist")
			return userNotExist
		}
		hash = make([]byte, len(stored_hash))
		copy(hash, stored_hash)

		return nil
	})

	if err != nil {
		switch err {
		case userNotExist:

			// create new key pair
			hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
			if err != nil {
				panic("Error hashing password")
			}

			// store username and hash
			err = db.Update(func(tx *bolt.Tx) error {
				bucket, err := tx.CreateBucketIfNotExists(bucket_users)
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
				panic("Error storing hash")
			}
			log.Println("Added new user")

		default:
			panic("Error opening Users bucket")
		}

	} else {

		// check if user is currently in a session
		err = db.Update(func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(bucket_tokens)
			if err != nil {
				panic("Error opening Tokens bucket")
			}
			token := bucket.Get(username)
			if token != nil {
				return userInSession
			}
			return nil
		})
		if err != nil {
			log.Println("User is currently in a session")
			return c, err // returning userInSession error
		}

		// compare stored hashed password and password from database
		err = bcrypt.CompareHashAndPassword(hash, password)
		if err != nil {
			log.Println("Incorrect password")
			return c, err // may want to return special error (incorrect password)
		}
	}

	log.Println("User is authenticated")

	// generate token
	token := make([]byte, 16)
	n, err := rand.Read(token)
	if n != 16 {
		panic("Error generating token")
	}

	// store token
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket_tokens)
		if err != nil {
			panic("Error opening Tokens bucket")
		}
		err = bucket.Put(username, token)
		if err != nil {
			panic("Error writing to Tokens bucket")
		}
		return nil
	})

	log.Println("Returning session token")

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
 * Handles an END request from the client
 */
func (s *SSLboardServer) End(_ context.Context, m *pb.Message) (*pb.Message, error) {

	log.Println("RPC call to End()")

	var bucket_tokens = []byte("Tokens")
	var userNotAuth = errors.New("User is not authenticated")
	var incorrectToken = errors.New("Incorrect session token")

	// open database
	db, err := bolt.Open("./board.db", 0666, nil)
	if err != nil {
		return m, err
	}
	defer db.Close()

	token := m.Token
	username := m.Username

	// verify that username has active token
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket_tokens)
		if bucket == nil {
			return userNotAuth
		}
		if token != string(bucket.Get([]byte(username))) {
			return incorrectToken
		}
		return nil
	})

	if err != nil {
		switch err {
		case userNotAuth:
			log.Println("User is currently in a session")
		case incorrectToken:
			log.Println("Incorrect session token")
		}
		return m, err
	}

	log.Println("User is logging out")

	// remove client token from list of active tokens
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(bucket_tokens)
		if err != nil {
			panic("Error opening Tokens bucket")
		}
		err = bucket.Delete([]byte(username))
		if err != nil {
			panic("Error writing to Tokens bucket")
		}
		return nil
	})

	log.Println("Active session terminated")

	return m, nil
}
