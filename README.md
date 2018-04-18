# SSLboard
Rutgers CS 419 OpenSSL Message Board Project

## Overview

We designed our project from the ground up, beginning with a simple client/server program, adding multithreading functionality, and slowly adding SSL certificates and end-to-end encryption.

We chose to use Go, a relatively new language that simplified a lot of C functionality, making it easier to implement a lot of the features required for this project. We used the `crypto/tls` package in Go for any SSL/TLS functions (certificates, TLS handshakes, etc.).

## Design

## Challenges

## Solutions

## Testing

We successfully tested the following events:

- [ ] Server can accept multiple clients
- [ ] Trying to fetch messages from a group that doesn’t exist
- [ ] Trying to provide an invalid username/password combo
- [ ] Attempting to submit blank messages or invalid group names
- [ ] Client cannot access private key, certificate, or password file

## Notes

We Used following commands for key and certificate generation

```sh
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
