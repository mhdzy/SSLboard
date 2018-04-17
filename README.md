# SSLboard
Rutgers CS 419 OpenSSL Message Board Project

## Overview

We designed our project from the ground up, beginning with a simple client/server program, adding multithreading functionality, and slowly adding SSL certificates and end-to-end encryption. 

We chose to use Go, a relatively new language that simplified a lot of C functionality, making it easier to implement a lot of the features required for this project.

Starting early definitely gave us a strong advantage, as did being able to freely choose a language to work with.

## Design

## Challenges

## Solutions



## Notes

Used following commands for key and certificate generation

```sh
openssl genrsa -out server.key 2048
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```
