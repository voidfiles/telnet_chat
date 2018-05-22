# Telnet Chat

[![CircleCI](https://circleci.com/gh/voidfiles/telnet_chat.svg?style=svg)](https://circleci.com/gh/voidfiles/telnet_chat)

A Simple chat server, a toy.

## Getting Started

This will install dependencies, including [dep](https://github.com/golang/dep) into your $GOPATH.  You'll be able to run tests and build the codebase afterwards.

```bash
make init
```

Run the tests

```bash
make test
```

To build

```bash
make build
```

## Run

After building you can run the command like this.

```
_work/telnet_chat_darwin_amd64 -configpath ./examples/config.toml
```

or you can use

```
make run
```

You can then telnet like this

```
telnet localhost 6000
```

You can also use the HTTP server with [httpie](https://httpie.org/).

```sh
$ http --verbose POST http://localhost:8000/api/messages client_id:=1 text="A message"

POST /api/messages HTTP/1.1
Accept: application/json, */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Content-Length: 37
Content-Type: application/json
Host: localhost:8000
User-Agent: HTTPie/0.9.9

{
    "client_id": 1,
    "text": "A message"
}

HTTP/1.1 200 OK
Content-Length: 79
Content-Type: application/json
Date: Tue, 22 May 2018 00:58:54 GMT

{
    "client_id": 1,
    "sent": "2018-05-21T18:58:54.226783057-06:00",
    "text": "A message"
}
```

```
$ http --verbose GET http://localhost:8000/api/messages

GET /api/messages HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate
Connection: keep-alive
Host: localhost:8000
User-Agent: HTTPie/0.9.9



HTTP/1.1 200 OK
Content-Length: 81
Content-Type: application/json
Date: Tue, 22 May 2018 00:59:44 GMT

[
    {
        "client_id": 1,
        "sent": "2018-05-21T18:58:54.226783057-06:00",
        "text": "A message"
    }
]
```

## Config

Here is an example config file.

```
telnetIp = "localhost"
telnetPort = "6000"
httpIp = "localhost"
httpPort = "8000"
logPath = "/tmp/logfile.json"
```

## Design

**ChatRoom**:
- Stores a history of messages
- Methods for adding messages and retrieving history
- Communicates via channels

**ClientPool**:
- Maintains client connections
- Distributes messages to all clients

**ChatServer**:
- Listens for inbound telnet connections
- Passes them to the ClientPool

**HttpServer**:
- Exposes an http interface to the ChatRoom

## Implmentation Details

To start I want to make sure I have all the basics covered.

* Config parsing
* Logging

Next I want to separate out the business logic and state of the chat server from the server. I can imagine a future choice to use something other than telnet as the connection protocol.

* ChatRoom
  - This should be in charge of the history of messages
  - This should accept and distribute messages to interested parties
* ClientPool
  - This should accept new network Connections
  - It should cleanup broken connections
  - It should listen for inbound messages from Clients and distribute them
  - It should listen for outbound messages and fan them out to Clients
* ChatServer
  - This should listen for new connections and send them to the ClientPool
* HttpServer
  - Should hook into methods provided by ChatRoom
  - To Expose message history
  - To allow
* Main
  - This should create an inbound message chan
  - This should create an outbound message chan
  - It should create a ClientPool and wire it up to the inbound and outbound chans
  - It should create a ChatRoom and wire it up to the inbound and outbound chans
  - It should create a ChatServer and wire it up to the ClientPool
  - It should then run the ClientPool, ChatRoom and ChatServer

## Considerations

If this were going to be something I would put into production I would consider a few things.

* Persistence
  - All the data is lost the second the server shuts down
  - If it was running on a single instance I might consider something like boltDB, or sqllite
* Bounds
  - I haven't fully vetted the project for bounds
  - I think their might be a few unbounded operations
    - Specifically the client fan out
* E2E tests
  - I think I could fully vet the service e2e
  - I'd want to build a system to do that
* HTTP Interface
  - I am not happy with my implementation
  - Frankly, I don't feel like I've ever nailed HTTP in go, but given a few iterations I think I could figure it out.
    - Mostly around patterns, good error handling, middleware, logging, stats, etc.
* Documentation
  - I feel like documentation should be used, and that use should guide its implementation. Every organization does it a little differently. In order to nail docs I would want to understand the goals of my organization, and then understand how documentation can best serve those goals.
* Stats
  - Again, different per org. I know in Go Prometheus is a common pattern, but I've also used statsd. I want to understand what the prevailing pattern is.

# Notes

Writing a telnet server is new for me. So, I used [golang-chat](https://github.com/kljensen/golang-chat) as inspiration for how to create a telnet server.
