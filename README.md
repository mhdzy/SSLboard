# SSLboard
Rutgers CS 419 OpenSSL Message Board Project

#### Matthew Handzy (mah376), Cullan Holloway (cph61), Viraj Patel (vcp25), Stephen Petrides (sdp157)

## Overview

We designed our project from the ground up, beginning with a simple client/server program, adding multithreading functionality, and slowly adding SSL certificates and end-to-end encryption.

We chose to use Go, a relatively new language that simplified a lot of C functionality, making it easier to implement a lot of the features required for this project. We used the `crypto/tls` package in Go for any SSL/TLS functions (certificates, TLS handshakes, etc.).

We also chose to use a fair amount of external libraries to 'trick out' our project, which we felt would ease a lot of development time -- and it did! Instead of spending time obsessing about delimiters, sending strings over the TLS pipe using Read/Write calls, and parsing on the server-side, gRPC calls and protobuffers allowed us to condense our message sending into structs over RPC calls. 

We had somewhat of a challenge figuring out boltDB, our database package. It was tricky and confusing setting this up, but offered a big learning experience in debugging since none of us had prior formal experience (academic or professional) with databases.

## Installation

In order to install Go, you must install the appropriate package from the following link:

    https://golang.org/dl/

Then, you need to place the source directory (.../SSLboard/) in $GOROOT/src/github.com/SleightOfHandzy/ ($GOROOT is typically $HOME/Go/) to end up with a relative pathname of: 

`~/Go/src/github.com/SleightOfHandzy/SSLboard/...`
   
Now, you will need to generate a `.key` and `.crt` file using the following commands: 

`openssl genrsa -out server.key 2048`

`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`
    
Next, you will need to place your `.key` and `.crt` files (openSSL generated key and certificate files) into the `/SSLboard/server/` directory. This is where our code assumes they will be, and it will compile but NOT work if they are not placed in this package. If you are using the compiled binaries in the `/bin/linux/` or `/bin/macOS/` folders (explained below), then please place your `.key` and `.crt` files directly into the respective folder.
    
From this point forward, you will be able to cd into a directory, compile, and execute our project. If you have trouble compiling our project in Windows, we advise you switch to a Linux or macOS machine.

## Execution

To be able to execute our project, you will need to first compile, then run both the server and client modules.

First, you will want to cd into the `/SSLboard/server/` module, and type: `go build .`

Then, you will need to run the server. Our code automatically runs it on port 8080, and this can be changed internally within the server.go file. You should have a go executable called `server` in your current directory. To execute this file, type: `./server`

Second, you will want to open a new Terminal window (or whatever shell you use) and cd into the /SSLboard/client/ module, and repeat the same build instruction. The '.' operator means that it will compile all .go files in the current directory. Just to re-iterate, you must type: `go build .`

Once this command finishes, you should have a go executable called `client` in your current directory. To execute this file, type: `./client server_IPaddress:8080`

Once these two are running in tandem and are connected, the client will ask you for authentication (username, password), and after authenticating yourself, you will be able to send commands to the Message Board.

The binaries for each OS are found in the `/SSLboard/bin/` directory. The subdirectories `/bin/macOS/` and `/bin/linux/` hold the compiled binaries for the respective operating systems, in the event that compiling using Go is difficult. The .crt and .key files should be present in each of the subdirectories.

## Packages

The external packages we used in this project include:

    google.golang.org/grpc
    google.golang.org/grpc/credentials
    golang.org/x/crypto/ssh/terminal
    golang.org/x/crypto/bcrypt
    github.com/golang/protobuf
    github.com/boltdb/bolt

We use `grpc` (Remote Procedure Calls) to simplify the interaction between the client and the server.

We use `grpc/credentials` to create transport credentials based on TLS. 

We use `crypto/ssh/terminal` to securely read passwords from the command line (not showing pwd characters == 'securely')

We use `crypto/bcrypt` for salting and hashing our passwords server-side.

We use `github.com/golang/protobuf` and a .proto file to implement our gRPC service, which is compiled into a .go file (found in SSLboard/pb) with the command 'protoc' (not included in this project).

We use `github.com/boltdb/bolt` as a back-end database, a key-value pair which uses 'buckets' for organization.

## Design

First, we implemented a simple client/server model in Go, using `TLS` sockets and `openSSL` to generate a `certificate` and a `key`. Then, we added `gRPC` functionality to simplify remote calls and eliminate parsing using delimiters (which could appear in a command, group name, message, and break our code). Once we had the basic structures set up, we added input via terminal: (username, password) verification, accepting a command (<CMD> <GRP> <MSG>), parsed that command, and sent a corresponding RPC to the correct server-side function in `/server/service.go`. 
  
We used the structs auto-generated by `protoc` (Message and Credentials structs) to send our complex messages over the pipe. You can directly view the result of the `/pb/SSLboard.proto` compilation in `/pb/SSLboard.pb.go`. This eliminated the need for delimitation and Read/Write calls over the TLS socket.

Then, we implemented our database to handle all calls for `GET` and `POST`. We store username-sessionTokenID as key-value pairs to validate that an RPC is being called by an authenticated (logged-in) user. If the database fails to find a match, the RPC was made by a non-authenticated user. 

We use a separate bucket for: username/password pairs, username/active_token pairs, and for each group name. If the named group's bucket does not exist, then the group does not exist. A no-message scenario is impossible, since each POST command requires a message as a parameter. We assume that the `/SSLboard/client/` and `/SSLboard/server/` modules do not know about each other, and therefore the client cannot read any files in the server's directory, such as the `board.db` database file or the `.key` and `.crt` files.

With `GET` calls, our service fetches all known messages in a certain group (bucket) and returns them to the client to quickly output to the terminal. If the <GRP> specified does not exist, an error will be returned.

With `POST` calls, our service opens a bucket to the <GRP> specified in the command arguments (error on non-existent bucket), and then appends the message to the end of the bucket. On success, a success message will be returned, and on a non-existent group name, an new bucket will be created for that group, and the message will be posted.
    
With `END` calls, our service removes the active username:token pairing from the database, effectively rendering the token associated with the user calling `END` inactive. Then the client will exit, and subsequent calls to the service using that un-authenticated token will be denied. The client will be able to log back in next time the client binary is executed.

#### gRPC

We setup a service, outlined in `/server/service.go`, to handle remote calls. In server.go, we have a main method that initiates a `gRPC service` and then calls `Serve()`, which allows the gRPC module to activate and handle requests. When an RPC is made, gRPC interprets the call and spawns a thread in the service to handle the function call (automatically). This simplifies multi-threading for us. Alternatively, we could have run a command `go funcName()`, which is Go's way to spawn a thread, however gRPC handles all of this for us. We have a very good understanding of multi-threading from the CS 214 and 416 experiences, so we felt comfortable abstracting this step of the process. Client-side, we open a connection to the corresponding gRPC server, and get handed back an object which allows us to make 'local calls' on the struct, which are interpreted by gRPC and handed to the service to execute.

#### User Session Tokens

Since we were using RPC's, theoretically, any connected (via `grpc.Dial()` call) user could make RPC's to our service. Thus, we needed to create a list of unique user session ID tokens, stored in the bolt database. We then authenticate each remote procedure call by checking the given username:token pair against our database, and then validate each call as coming from an authenticated user.

## Challenges

#### Learning Golang (Go)

To accomplish this project, our entire team needed to learn Go from scratch. None of our members had any real prior exposure or experience with this language, however we started early and this allowed us to explore many options and truly see the power of the Go library and external packages.

#### User Session Tokens

As explained above, user session tokens protect foreign clients from using RPC's to access our database. This step in development took a while, notably because we had a design choice that we would not allow multi-client logins (same user/pass combination on different clients). Thus, we had to ensure that only one instance of a username was logged in at any given time. Unfortunately, this step took relatively the longest compared to how much time we expected to spend on it.

#### gRPC

Learning gRPC was a bit of a challenge, especially figuring out how we would implement the compiled .go file (compiled from .proto file). However, after a fair bit of research and testing, we were able to successfully implement RPC's which made sending information (credentials, messages, anything) over the TLS pipe trivial.

#### protobuffers

Learning how to write a `.proto` file was rather easy, as we only had to define a few methods that we were to implement along with the necessary structs to send information over our connection. However, figuring out how to incorporate gRPC's functionality along with implementing our interfaces and structs took a bit of time. 

#### boltDB

BoltDB had a very high learning curve for our team, as none of us had any exposure to this particular database before. Understanding the transaction syntax was a bit difficult, and handling race conditions that were created during development. Overall, this part of the project helped us simplify reading/writing to persistent memory (storing the messageBoard), but took some time to really understand and implement.

#### dep (ensure)

dep allowed us to bundle all external dependencies (the not-standard Go packages) and export them with our project. This offers an extremely large convinience to us, and to whoever is running our project by not requiring the executor to `go get` (command-line tool for downloading Go packages) any external packages. All of our dependencies can be found in the `/SSLboard/vendor/` directory.

## Solutions

We powered through all of these challenges to present a final product that our team is very proud of. We learned a lot while developing this project and are now better programmers and have a better feel for Golang, gRPC, Protocol Buffers, and boltDB.

## Testing

We successfully tested the following mandatory events:

- [x] Server can accept multiple clients
- [x] Trying to fetch messages from a group that doesn’t exist
- [x] Trying to provide an invalid username/password combo
- [x] Attempting to submit blank messages or invalid group names
- [x] Client cannot access private key, certificate, or password file

We also successfully tested for our own purposes:

- [x] Client program will inform user on incorrect command syntax (bad commands are handled in the way that they are not sent to the server, which avoids missing paramters in RPC calls)
- [x] Database errors/problems: buckets not existing, buckets being empty, race conditions
- [x] Multi-session logins disabled (most-current login is the considered the 'real' one)
- [x] Trying to make an RPC without a server issued token, a previously useable token, or a valid token that is issued for a different user (invalid user/token combo)
- [x] END successfully removes a client token so that a user must re-authenticate to use the formerly active token.
- [x] When the server restarts, active tokens are flushed and a new session must be authenticated for each formerly connected client.
