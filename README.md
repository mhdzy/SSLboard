# SSLboard
Rutgers CS 419 OpenSSL Message Board Project

## Overview

We designed our project from the ground up, beginning with a simple client/server program, adding multithreading functionality, and slowly adding SSL certificates and end-to-end encryption.

We chose to use Go, a relatively new language that simplified a lot of C functionality, making it easier to implement a lot of the features required for this project. We used the `crypto/tls` package in Go for any SSL/TLS functions (certificates, TLS handshakes, etc.).

## Installation and Packages

In order to install Go, you must install the appropriate package from the following link:

```sh
https://golang.org/dl/
```

Then, you need to run this command (to install a necessary std. library):

```sh
go get golang.org/x/crypto/ssh/terminal
```

grpc:
```sh
go get -u google.golang.org/grpc
```

protoc plugin:
```sh
go get -u github.com/golang/protobuf/protoc-gen-go
```

bcrypt plugin:
```sh
go get -u golang.org/x/crypto/bcrypt
```

protoc:
```sh
brew install protoc
```

## Design

## Challenges

## Solutions

## Testing

We successfully tested the following events:

- [x] Server can accept multiple clients
- [ ] Trying to fetch messages from a group that doesn’t exist
- [ ] Trying to provide an invalid username/password combo
- [ ] Attempting to submit blank messages or invalid group names
- [ ] Client cannot access private key, certificate, or password file

## Notes

We used following commands for key and certificate generation

```sh
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
